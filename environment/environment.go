package environment

import (
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/container"
	"log"
)

type Environment struct {
	Client    *container.Client // The docker client
	BaseImage string            // The docker image
	Running   []string          // IDs of containers currently running
	Timeout   int               // In seconds
	Logger    *log.Logger
	Commands  []string
}

// Create a new Environment instance.
func NewEnvironment(baseImage string, commands []string, client *container.Client, logger *log.Logger) *Environment {
	return &Environment{
		Client:    client,
		BaseImage: baseImage,
		Running:   []string{},
		Logger:    logger,
		Commands:  commands,
	}
}

func (env *Environment) Run() error {
	resp, err := env.Client.CreateContainer(env.BaseImage, env.Commands)
	if err != nil {
		return err
	}
	env.Running = append(env.Running, resp.ID)
	_, err = env.Client.WaitForContainer(resp.ID, env.Timeout)
	if err != nil {
		return err
	}
	env.Logger.Printf(fmt.Sprintf("Container [%s] finished!", resp.ID))
	return nil
}
