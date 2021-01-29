// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutils

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Container is a combination builder and testcontainers.Container wrapper
// for convenient creation and management of docker images and containers.
type Container struct {
	Image        string
	Dockerfile   testcontainers.FromDockerfile
	Env          map[string]string
	ExposedPorts []string
	WaitingFor   []wait.Strategy
	container    *testcontainers.Container
	req          testcontainers.ContainerRequest
}

var _ testcontainers.Container = (*Container)(nil)

// To be used as a builder whose BuildRequest() method provides the actual instance capable of being started, and who
// implements a testcontainers.Container.
func NewContainer() Container {
	return Container{
		Env: map[string]string{},
	}
}

func (container Container) WithImage(image string) Container {
	container.Image = image
	return container
}

func (container Container) WithDockerfile(dockerfile string) Container {
	container.Dockerfile.Dockerfile = dockerfile
	return container
}

func (container Container) WithContext(path string) Container {
	container.Dockerfile.Context = path
	return container
}

func (container Container) WithContextArchive(contextArchive io.Reader) Container {
	container.Dockerfile.ContextArchive = contextArchive
	return container
}

func (container Container) WithEnv(env map[string]string) Container {
	for k, v := range env {
		container.Env[k] = v
	}
	return container
}

func (container Container) WithEnvVar(key, value string) Container {
	container.Env[key] = value
	return container
}

func (container Container) WithExposedPorts(ports []string) Container {
	container.ExposedPorts = append(container.ExposedPorts, ports...)
	return container
}

func (container Container) WithExposedPort(port string) Container {
	container.ExposedPorts = append(container.ExposedPorts, port)
	return container
}

func (container Container) WillWaitForPort(port string) Container {
	container.WaitingFor = append(container.WaitingFor, wait.ForListeningPort(nat.Port(port)))
	return container
}

func (container Container) WillWaitForLog(logStatement string) Container {
	container.WaitingFor = append(container.WaitingFor, wait.ForLog(logStatement))
	return container
}

// Acts as the build method.
func (container Container) BuildRequest() *Container {
	container.req = testcontainers.ContainerRequest{
		Image:          container.Image,
		FromDockerfile: container.Dockerfile,
		Env:            container.Env,
		ExposedPorts:   container.ExposedPorts,
		WaitingFor:     wait.ForAll(container.WaitingFor...),
	}
	return &container
}

func (container *Container) Start(ctx context.Context) error {
	started, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: container.req,
		Started:          true,
	})
	container.container = &started
	return err
}

func (container *Container) assertStarted(operation string) error {
	if container.container == nil || (*container.container) == nil {
		return fmt.Errorf("cannot %s() unstarted container", operation)
	}
	return nil
}

func (container *Container) Stop(ctx context.Context) error {
	if err := container.assertStarted("Stop"); err != nil {
		return err
	}
	return container.Terminate(ctx)
}

func (container *Container) GetContainerID() string {
	if err := container.assertStarted("GetContainerID"); err != nil {
		return ""
	}
	return (*container.container).GetContainerID()
}

func (container *Container) Endpoint(ctx context.Context, s string) (string, error) {
	if err := container.assertStarted("Endpoint"); err != nil {
		return "", err
	}
	return (*container.container).Endpoint(ctx, s)
}

func (container *Container) PortEndpoint(ctx context.Context, port nat.Port, s string) (string, error) {
	if err := container.assertStarted("PortEndpoint"); err != nil {
		return "", err
	}
	return (*container.container).PortEndpoint(ctx, port, s)
}

func (container *Container) Host(ctx context.Context) (string, error) {
	if err := container.assertStarted("Host"); err != nil {
		return "", err
	}
	return (*container.container).Host(ctx)
}

func (container *Container) MappedPort(ctx context.Context, port nat.Port) (nat.Port, error) {
	if err := container.assertStarted("MappedPort"); err != nil {
		return "", err
	}
	return (*container.container).MappedPort(ctx, port)
}

func (container *Container) Ports(ctx context.Context) (nat.PortMap, error) {
	if err := container.assertStarted("Ports"); err != nil {
		return nil, err
	}
	return (*container.container).Ports(ctx)
}

func (container *Container) SessionID() string {
	if err := container.assertStarted("SessionID"); err != nil {
		return ""
	}
	return (*container.container).SessionID()
}

func (container *Container) Terminate(ctx context.Context) error {
	if err := container.assertStarted("Terminate"); err != nil {
		return err
	}
	return (*container.container).Terminate(ctx)
}

func (container *Container) Logs(ctx context.Context) (io.ReadCloser, error) {
	if err := container.assertStarted("Logs"); err != nil {
		return nil, err
	}
	return (*container.container).Logs(ctx)
}

func (container *Container) FollowOutput(consumer testcontainers.LogConsumer) {
	if err := container.assertStarted("FollowOutput"); err == nil {
		(*container.container).FollowOutput(consumer)
	}
}

func (container *Container) StartLogProducer(ctx context.Context) error {
	if err := container.assertStarted("StartLogProducer"); err != nil {
		return err
	}
	return (*container.container).StartLogProducer(ctx)
}

func (container *Container) StopLogProducer() error {
	if err := container.assertStarted("StopLogProducer"); err != nil {
		return err
	}
	return (*container.container).StopLogProducer()
}

func (container *Container) Name(ctx context.Context) (string, error) {
	if err := container.assertStarted("Name"); err != nil {
		return "", err
	}
	return (*container.container).Name(ctx)
}

func (container *Container) Networks(ctx context.Context) ([]string, error) {
	if err := container.assertStarted("Network"); err != nil {
		return nil, err
	}
	return (*container.container).Networks(ctx)
}

func (container *Container) NetworkAliases(ctx context.Context) (map[string][]string, error) {
	if err := container.assertStarted("NetworkAliases"); err != nil {
		return nil, err
	}
	return (*container.container).NetworkAliases(ctx)
}

func (container *Container) Exec(ctx context.Context, cmd []string) (int, error) {
	if err := container.assertStarted("Exec"); err != nil {
		return 0, err
	}
	return (*container.container).Exec(ctx, cmd)
}

func (container *Container) ContainerIP(ctx context.Context) (string, error) {
	if err := container.assertStarted("ContainerIP"); err != nil {
		return "", err
	}
	return (*container.container).ContainerIP(ctx)
}

func (container *Container) CopyFileToContainer(ctx context.Context, hostFilePath string, containerFilePath string, fileMode int64) error {
	if err := container.assertStarted("CopyFileToContainer"); err != nil {
		return err
	}
	return (*container.container).CopyFileToContainer(ctx, hostFilePath, containerFilePath, fileMode)
}