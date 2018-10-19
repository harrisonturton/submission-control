package environment

import (
	"bytes"
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/container"
	"log"
	"time"
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

// RunWithLogs creates & runs a new container.
func (env *Environment) Run() (*string, error) {
	resp, err := env.Client.CreateContainer(env.BaseImage, env.Commands)
	if err != nil {
		return nil, err
	}
	env.Running = append(env.Running, resp.ID)
	if err = env.Client.WaitForContainer(resp.ID, time.Second*10); err != nil {
		return nil, err
	}
	env.Logger.Printf(fmt.Sprintf("Container exited: %s", resp.ID))
	return &resp.ID, nil
}

// RunWithLogs creates & runs a new container, but also optionally returns the STDOUT
// and STDERR.
func (env *Environment) RunWithLogs(showStdout bool, showStderr bool) (string, error) {
	id, err := env.Run()
	if err != nil {
		return "", err
	}
	logReader, err := env.Client.ReadContainerLogs(*id, showStdout, showStderr)
	if err != nil {
		return "", err
	}
	defer logReader.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(logReader)
	return buf.String(), nil
}
