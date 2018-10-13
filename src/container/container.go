package container

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Client struct {
	instance *client.Client
	context  context.Context
}

// Create a new instance of a Docker client.
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
func (client *Client) CreateContainer(fromImageID string) (container.ContainerCreateCreatedBody, error) {
	return client.instance.ContainerCreate(client.context, &container.Config{
		Image: fromImageID,
	}, nil, nil, "")
}

// Start a container. Can start a stopped container, or a container that hasn't
// been run yet.
func (client *Client) StartContainer(ID string) error {
	return client.instance.ContainerStart(client.context, ID, types.ContainerStartOptions{})
}

// Check if the files in the container have changed since it has been run.
// After we run & clean up user code, check this to make sure no side-effects have
// occured. If they have, restart the container and flag the user's submission.
func (client *Client) ContainerEnvHasChanged(containerID string) (bool, error) {
	changes, err := client.instance.ContainerDiff(client.context, containerID)
	if err != nil {
		return false, err
	}
	return len(changes) > 0, nil
}
