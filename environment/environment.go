package environment

import (
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/container"
)

type Environment struct {
	Client    *container.Client // The docker client
	BaseImage string            // The docker image
	Running   []string          // IDs of containers currently running
	Timeout   int               // In seconds
}

// Create a new Environment instance.
func NewEnvironment(baseImage string, client *container.Client) *Environment {
	return &Environment{
		Client:    client,
		BaseImage: baseImage,
		Running:   []string{},
	}
}

func (env *Environment) Run() error {
	resp, err := env.Client.CreateContainer(env.BaseImage)
	if err != nil {
		return err
	}
	env.Running = append(env.Running, resp.ID)
	_, err = env.Client.WaitForContainer(resp.ID, env.Timeout)
	if err != nil {
		return err
	}
	fmt.Println("Container finished! " + resp.ID)
	return nil
}
