package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"io"
	"time"
)

// Client is an interface to the Docker daemon. It abstracts all the library
// calls under one unifying interface.
type Client struct {
	Instance *client.Client
	Context  context.Context
}

// NewClient creates a new interface to the Docker daemon. If the GO SDK version
// is incompatible with the API version, try specifying a lower version number.
func NewClient(version string) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		return nil, err
	}
	return &Client{cli, context.Background()}, nil
}

// CreateContainer creates a new container from an existing image. Note, this does
// NOT run the container.
func (client *Client) CreateContainer(fromImageID string, commands []string) (container.ContainerCreateCreatedBody, error) {
	return client.Instance.ContainerCreate(client.Context, &container.Config{
		Image: fromImageID,
		Cmd:   commands,
	}, nil, nil, "")
}

// StartContainer launches a non-running container.
func (client *Client) StartContainer(containerID string) error {
	return client.Instance.ContainerStart(
		client.Context, containerID, types.ContainerStartOptions{})
}

// ReadContainerLogs returns the Stdout or Stderr (or both) of a container.
func (client *Client) ReadContainerLogs(containerID string, showStdout bool, showStderr bool) (io.ReadCloser, error) {
	return client.Instance.ContainerLogs(
		client.Context, containerID, types.ContainerLogsOptions{
			ShowStdout: showStdout,
			ShowStderr: showStderr,
		})
}

// WaitForContainer will start a container, and block until it either finishes,
// or the runtime surpasses the timeout.
func (client *Client) WaitForContainer(containerID string, timeout time.Duration) error {
	if err := client.StartContainer(containerID); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(client.Context, timeout)
	defer cancel()
	respCh, errCh := client.Instance.ContainerWait(client.Context, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return err
	case <-respCh:
		return nil
	case <-ctx.Done(): // Timeout
		return errors.New(fmt.Sprintf("WaitForContainer timeout exceeded: ", containerID))
	}
}

// CreateService creates a new service from an existing image
func (client *Client) CreateService(fromImageID string, replicas uint64, command []string) (types.ServiceCreateResponse, error) {
	return client.Instance.ServiceCreate(client.Context, swarm.ServiceSpec{
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image:   fromImageID,
				Command: command,
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			},
		},
	}, types.ServiceCreateOptions{})
}

// RemoveService stops the replicas of an en existing service, and removes it from the
// docker daemon.
func (client *Client) RemoveService(serviceID string) error {
	return client.Instance.ServiceRemove(client.Context, serviceID)
}

// ScaleService changes the number of replicas on a service.
func (client *Client) ScaleService(serviceID string, replicas uint64) error {
	// Need to make sure *Spec and Version numbers match
	inspectResp, _, err := client.Instance.ServiceInspectWithRaw(client.Context, serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}
	spec := inspectResp.Spec
	spec.Mode.Replicated.Replicas = &replicas
	_, err = client.Instance.ServiceUpdate(client.Context, serviceID, inspectResp.Meta.Version, spec, types.ServiceUpdateOptions{})
	return err
}
