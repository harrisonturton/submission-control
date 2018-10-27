package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"time"
)

// Client is the interface to a container management
// platform
type Client interface {
	Create(imageID string) (string, error)
	Start(containerID string) error
	ReadLogs(containerID string, showStdout bool, showStderr bool) (string, error)
	Wait(containerID string, timeout time.Duration) error
}

// Docker is a connection to the Docker daemon.
type Docker struct {
	Instance *client.Client
	Context  context.Context
}

// New creates a new interface to the Docker daemon. If the GO SDK version
// is incompatible with the API version, try specifying a lower version number.
func New(version string) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		return nil, err
	}
	return &Docker{cli, context.Background()}, nil
}

// Create creates a new container from an existing image. Note, this does
// NOT run the container.
func (client *Docker) Create(imageID string) (string, error) {
	resp, err := client.Instance.ContainerCreate(client.Context, &container.Config{
		Image: imageID,
	}, nil, nil, "")
	return resp.ID, err
}

// Start launches a non-running container.
func (client *Docker) Start(containerID string) error {
	return client.Instance.ContainerStart(
		client.Context, containerID, types.ContainerStartOptions{})
}

// ReadLogs returns the Stdout or Stderr (or both) of a container.
func (client *Docker) ReadLogs(containerID string, showStdout bool, showStderr bool) (string, error) {
	logReader, err := client.Instance.ContainerLogs(
		client.Context, containerID, types.ContainerLogsOptions{
			ShowStdout: showStdout,
			ShowStderr: showStderr,
		})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(logReader)
	logReader.Close()
	return buf.String(), nil
}

// Wait will start a container, and block until it either finishes,
// or the runtime surpasses the timeout.
func (client *Docker) Wait(containerID string, timeout time.Duration) error {
	if err := client.Start(containerID); err != nil {
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
		return fmt.Errorf("WaitForContainer timeout exceeded: ", containerID)
	}
}

// CreateService creates a new service from an existing image
func (client *Docker) CreateService(fromImageID string, replicas uint64, command []string) (types.ServiceCreateResponse, error) {
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
func (client *Docker) RemoveService(serviceID string) error {
	return client.Instance.ServiceRemove(client.Context, serviceID)
}

// ScaleService changes the number of replicas on a service.
func (client *Docker) ScaleService(serviceID string, replicas uint64) error {
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