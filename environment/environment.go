package environment

import (
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/container"
	"io"
	"io/ioutil"
	"log"
	"os"
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

func (env *Environment) Run() (*string, error) {
	resp, err := env.Client.CreateContainer(env.BaseImage, env.Commands)
	if err != nil {
		return nil, err
	}
	env.Running = append(env.Running, resp.ID)
	_, err = env.Client.WaitForContainer(resp.ID, env.Timeout)
	if err != nil {
		return nil, err
	}
	env.Logger.Printf(fmt.Sprintf("Container exited: %s", resp.ID))
	return &resp.ID, nil
}

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
	io.Copy(os.Stdout, logReader)
	data, err := ioutil.ReadAll(logReader)
	if err != nil && err != io.EOF {
		return "", err
	}
	fmt.Printf("Data: %s\n", string(data))
	return string(data), nil
}
