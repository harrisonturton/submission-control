package container

import (
	"fmt"
	"github.com/harrisonturton/hydra-ci/util"
	"time"
)

type Container struct {
	ID            string
	EnvironmentID string
	BaseImage     string
	Results       chan<- string
}

// Create a new container instance
func NewContainer(environmentID string, baseImage string, results chan<- string) *Container {
	return &Container{
		ID:            util.Uuid(),
		EnvironmentID: environmentID,
		BaseImage:     baseImage,
		Results:       results,
	}
}

// Fake building & running test cases
func (container *Container) Run(request string) {
	fmt.Println(fmt.Sprintf("[%s][%s] Handling %s", container.EnvironmentID, container.ID, request))
	time.Sleep(time.Second * 5)
	container.Results <- fmt.Sprintf("[%s][%s] Passed [%s]", container.EnvironmentID, container.ID, request)
}
