package remote

import (
	"github.com/harrisonturton/submission-control/daemon/client"
	"github.com/harrisonturton/submission-control/daemon/types"
	"log"
)

type RemoteServer struct {
	client *client.Client
	logger *log.Logger
}

// Create a new remote server attached to a daemon instance.
func NewRemoteServer(client *client.Client, logger *log.Logger) *RemoteServer {
	return &RemoteServer{
		client: client,
		logger: logger,
	}
}

// Handle an RPC call to create a service
func (remote *RemoteServer) CreateService(createArgs *types.ServiceCreateSpec, ID *string) error {
	remote.logger.Printf("Recieved RPC call to CreateService with args: %s", createArgs)
	resp, err := remote.client.CreateService(createArgs.BaseImage, createArgs.Replicas, createArgs.Commands)
	if err != nil {
		return err
	}
	*ID = resp.ID
	return err
}

// Handle RPC call to remove a service
func (remote *RemoteServer) RemoveService(removeArgs *types.ServiceRemoveSpec, reply *int) error {
	remote.logger.Printf("Recieved RPC call to RemoveService with args: %s", removeArgs)
	return remote.client.RemoveService(removeArgs.ServiceID)
}

// Handle RPC call to scale a service (i.e. increase or decrease number of replicas)
func (remote *RemoteServer) ScaleService(scaleArgs *types.ServiceScaleSpec, reply *int) error {
	remote.logger.Printf("Recieved RPC call to ScaleService with args: %s", scaleArgs)
	return remote.client.ScaleService(scaleArgs.ServiceID, scaleArgs.Replicas)
}
