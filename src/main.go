package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Manager struct {
	Cli *client.Client
	Ctx context.Context
}

func NewManager(version string) (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		return nil, err
	}
	return &Manager{cli, context.Background()}, nil
}

// Create a new container from an existing image.
func (client *Manager) CreateContainer(fromImageID string) (container.ContainerCreateCreatedBody, error) {
	return client.Cli.ContainerCreate(client.Ctx, &container.Config{
		Image: fromImageID,
	}, nil, nil, "")
}

// Want to try:
// 1. Start new docker image (submission base) DONE!
// 2. Copy project into it (submission content)
// 3. Run it
// 4. Collect results

func main() {
	manager, err := NewManager("1.38")
	panicErr(err)

	resp, err := manager.CreateContainer("ada-project")
	panicErr(err)

	statusCh, errCh := manager.Cli.ContainerWait(manager.Ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		panicErr(err)
	case <-statusCh:
	}

	out, err := manager.Cli.ContainerLogs(manager.Ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	panicErr(err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	s := buf.String()
	fmt.Println(s)
}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
