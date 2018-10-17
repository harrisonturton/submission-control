package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"time"
)

type Client struct {
	Instance *client.Client
	Context  context.Context
}

// Create a new Instance of a Docker client.
// If the Go SDK version is incompatible with the
// Docker client service, specify a lower version no.
func NewClient(version string) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		return nil, err
	}
	return &Client{cli, context.Background()}, nil
}

// Create a new container from an existing image.
func (client *Client) CreateContainer(fromImageID string, commands []string) (container.ContainerCreateCreatedBody, error) {
	return client.Instance.ContainerCreate(client.Context, &container.Config{
		Image: fromImageID,
		Cmd:   commands,
	}, nil, nil, "")
}

// Wait for the next exit state of a container, with timeout (in seconds)
func (client *Client) WaitForContainer(containerID string, timeout int) (*container.ContainerWaitOKBody, error) {
	ctx, cancel := context.WithTimeout(client.Context, time.Second*10)
	defer cancel()
	respCh, errCh := client.Instance.ContainerWait(client.Context, containerID, container.WaitConditionNextExit)
	select {
	case err := <-errCh:
		return nil, err
	case resp := <-respCh:
		return &resp, nil
	case <-ctx.Done(): // Timeout
		return nil, errors.New(fmt.Sprintf("WaitForContainer timeout exceeded for container %s", containerID))
	}
}

// Create a new service from an exiting image
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

// Remove an existing service
func (client *Client) RemoveService(serviceID string) error {
	return client.Instance.ServiceRemove(client.Context, serviceID)
}

// Change the number of replicas on a service
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

// Start a container. Can start a stopped container, or a container that hasn't
// been run yet.
func (client *Client) StartContainer(ID string) error {
	return client.Instance.ContainerStart(client.Context, ID, types.ContainerStartOptions{})
}

// Check if the files in the container have changed since it has been run.
// After we run & clean up user code, check this to make sure no side-effects have
// occured. If they have, restart the container and flag the user's submission.
func (client *Client) ContainerEnvHasChanged(containerID string) (bool, error) {
	changes, err := client.Instance.ContainerDiff(client.Context, containerID)
	if err != nil {
		return false, err
	}
	return len(changes) > 0, nil
}
