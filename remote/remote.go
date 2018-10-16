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

// Handle and RPC call to container.ScaleService
func (remote *RemoteServer) ScaleContainer(scaleArgs *types.ScaleArgs, reply *int) error {
	return remote.client.ScaleService(scaleArgs.ServiceID, scaleArgs.Replicas)
}
