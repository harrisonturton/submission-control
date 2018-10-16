package remote

import (
	"github.com/harrisonturton/hydra-daemon/container"
	"github.com/harrisonturton/hydra-daemon/types"
)

type RemoteServer struct {
	client *container.Client
}

// Create a new remote server instance
func NewRemoteServer(version string) (*RemoteServer, error) {
	client, err := container.NewClient(version)
	if err != nil {
		return nil, err
	}
	return &RemoteServer{client}, nil
}

// Handle an RPC call to create a service
func (remote *RemoteServer) CreateService(createArgs *types.ServiceCreateSpec, ID *string) error {
	resp, err := remote.client.CreateService(createArgs.BaseImage, createArgs.Replicas, createArgs.Commands)
	if err != nil {
		return err
	}
	*ID = resp.ID
	return err
}

// Handle RPC call to remove a service
func (remote *RemoteServer) RemoveService(removeArgs *types.ServiceRemoveSpec, reply *int) error {
	return remote.client.RemoveService(removeArgs.ServiceID)
}

// Handle RPC call to scale a service (i.e. increase or decrease number of replicas)
func (remote *RemoteServer) ScaleService(scaleArgs *types.ServiceScaleSpec, reply *int) error {
	return remote.client.ScaleService(scaleArgs.ServiceID, scaleArgs.Replicas)
}
