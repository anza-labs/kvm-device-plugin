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

package plugin

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func DevicePluginServer(plugin v1beta1.DevicePluginServer) *grpc.Server {
	srv := grpc.NewServer()
	v1beta1.RegisterDevicePluginServer(srv, plugin)
	return srv
}

func RegisterDevicePlugin(ctx context.Context, name, socket string) error {
	conn, err := grpc.NewClient(
		filepath.Join(v1beta1.KubeletSocket),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", addr)
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to kubelet: %v", err)
	}
	defer conn.Close() //nolint:errcheck // best effort call

	_, err = v1beta1.NewRegistrationClient(conn).Register(ctx, &v1beta1.RegisterRequest{
		Version:      v1beta1.Version,
		ResourceName: name,
		Endpoint:     filepath.Base(socket),
	})
	if err != nil {
		return fmt.Errorf("failed to register plugin with kubelet service: %v", err)
	}

	return nil
}
