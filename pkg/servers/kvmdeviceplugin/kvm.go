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

package kvmdeviceplugin

import (
	"context"

	"github.com/anza-labs/kvm-device-plugin/pkg/plugin"

	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Server struct {
	name   string
	socket string
}

var _ v1beta1.DevicePluginServer = (*Server)(nil)

func New(name, socket string) *Server {
	return &Server{
		name:   name,
		socket: socket,
	}
}

func (s *Server) Register(ctx context.Context) error {
	// TODO: wait for Server to be running

	return plugin.RegisterDevicePlugin(ctx, s.name, s.socket)
}

// Manager.
func (s *Server) GetDevicePluginOptions(
	ctx context.Context,
	_ *v1beta1.Empty,
) (*v1beta1.DevicePluginOptions, error) {
	return &v1beta1.DevicePluginOptions{
		PreStartRequired:                false,
		GetPreferredAllocationAvailable: false,
	}, nil
}

// returns the new list.
func (s *Server) ListAndWatch(
	_ *v1beta1.Empty,
	lws v1beta1.DevicePlugin_ListAndWatchServer,
) error {
	panic("unimplemented")
}

// GetPreferredAllocation returns a preferred set of devices to allocate
// from a list of available ones. The resulting preferred allocation is not
// guaranteed to be the allocation ultimately performed by the
// devicemanager. It is only designed to help the devicemanager make a more
// informed allocation decision when possible.
func (s *Server) GetPreferredAllocation(
	ctx context.Context,
	req *v1beta1.PreferredAllocationRequest,
) (*v1beta1.PreferredAllocationResponse, error) {
	panic("unimplemented")
}

// of the steps to make the Device available in the container.
func (s *Server) Allocate(
	ctx context.Context,
	req *v1beta1.AllocateRequest,
) (*v1beta1.AllocateResponse, error) {
	panic("unimplemented")
}

// such as resetting the device before making devices available to the container.
func (s *Server) PreStartContainer(
	ctx context.Context,
	req *v1beta1.PreStartContainerRequest,
) (*v1beta1.PreStartContainerResponse, error) {
	panic("unimplemented")
}
