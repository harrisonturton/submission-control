package environment

import (
	"fmt"
	"github.com/harrisonturton/hydra-ci/container"
	"github.com/harrisonturton/hydra-ci/util"
	"sync"
)

type Environment struct {
	ID       string
	Base     string                    // The base container image
	Requests <-chan string             // Where we recieve testing requests
	Results  chan<- string             // Where we send testing results
	Free     chan *container.Container // Free containers
	Count    int                       // Number of containers in action
}

const (
	IdleContainers = 2 // Minimum number of containers
	MaxContainers  = 5 // Maxmimum number of containers
)

// Create a new Environment instance
func NewEnvironment(base string, requests <-chan string, results chan<- string) *Environment {
	// Initialize idle containers
	id := util.Uuid()
	free := make(chan *container.Container, MaxContainers)
	for i := 0; i < IdleContainers; i++ {
		free <- container.NewContainer(id, "comp2310-assignment-1", results)
	}
	return &Environment{
		ID:       id,
		Base:     base,
		Requests: requests,
		Results:  results,
		Free:     free,
		Count:    IdleContainers,
	}
}

// Main loop. Spin up a new container for every request,
// until the maximum number of containers is reached.
func (env *Environment) Run(wg *sync.WaitGroup, stop <-chan bool) {
	for {
		select {
		case request := <-env.Requests:
			fmt.Println(fmt.Sprintf("[%s] Handling %s", env.ID, request))
			go env.HandleRequest(request)
		case <-stop:
			defer wg.Done()
			return
		}
	}
}

// Either use a spare container, or spin up a new one.
func (env *Environment) HandleRequest(request string) {
	select {
	// Use a free container if we have one
	case cont := <-env.Free:
		cont.Run(request)
		env.Free <- cont
	// Otherwise create a new container if we can
	default:
		if env.Count < MaxContainers {
			env.Count++
			cont := container.NewContainer(env.ID, "comp2310-assignment-1", env.Results)
			cont.Run(request)
			env.Free <- cont
		} else {
			cont := <-env.Free
			cont.Run(request)
			env.Free <- cont
		}
	}
}
