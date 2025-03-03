// Copyright 2025 anza-labs contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/anza-labs/kvm-device-plugin/pkg/plugin"
	"github.com/anza-labs/kvm-device-plugin/pkg/servers/kvmdeviceplugin"
)

const (
	pluginName = "virt.io/kvm"

	gracePeriod = 5 * time.Second
)

func main() {
	// TODO: any logger settup must be done here, before first log call.
	log := slog.Default()

	if err := run(context.Background(), log, mainOptions{
		pluginEndpoint: "unix:///tmp/test.sock",
	}); err != nil {
		slog.Error("Critical failure", "error", err)
		os.Exit(1)
	}
}

type mainOptions struct {
	pluginEndpoint string
}

func run(ctx context.Context, log *slog.Logger, opts mainOptions) error {
	ctx, stop := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	log.InfoContext(ctx, "Starting plugin")
	eg, ctx := errgroup.WithContext(ctx)
	kvm := kvmdeviceplugin.New(pluginName, opts.pluginEndpoint)
	server := plugin.DevicePluginServer(kvm)

	eg.Go(func() error {
		return shutdown(ctx, server)
	})
	eg.Go(func() error {
		return plugin.RegisterDevicePlugin(ctx, pluginName, opts.pluginEndpoint)
	})
	eg.Go(func() error {
		lis, cleanup, err := listener(ctx, opts.pluginEndpoint)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}
		defer cleanup()
		return server.Serve(lis)
	})

	return eg.Wait()
}

func listener(
	ctx context.Context,
	pluginEndpoint string,
) (net.Listener, func(), error) {
	endpointURL, err := url.Parse(pluginEndpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse COSI endpoint: %w", err)
	}

	listenConfig := net.ListenConfig{}

	if endpointURL.Scheme == "unix" {
		// best effort call to remove the socket if it exists, fixes issue with restarted pod that did not exit gracefully
		_ = os.Remove(endpointURL.Path)
	}

	listener, err := listenConfig.Listen(ctx, endpointURL.Scheme, endpointURL.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create listener: %w", err)
	}

	cleanup := func() {
		if err := listener.Close(); err != nil {
			slog.Error("Failed to close listener", "error", err)
		}

		if endpointURL.Scheme == "unix" {
			if err := os.Remove(endpointURL.Path); err != nil {
				slog.Error("Failed to remove old socket", "error", err)
			}
		}
	}

	return listener, cleanup, nil
}

func shutdown(ctx context.Context, server *grpc.Server) error {
	<-ctx.Done()
	slog.Info("Shutting down")
	dctx, stop := context.WithTimeout(context.Background(), gracePeriod)
	defer stop()

	c := make(chan struct{})
	if server != nil {
		go func() {
			server.GracefulStop()
			c <- struct{}{}
		}()
		for {
			select {
			case <-dctx.Done():
				slog.Info("Forcing shutdown")
				server.Stop()
				return nil
			case <-c:
				return nil
			}
		}
	}
	return nil
}
