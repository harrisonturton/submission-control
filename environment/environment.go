package environment

import (
	"fmt"
	"sync"
)

type Environment struct {
	ID       string
	Requests <-chan string
	Results  chan<- string
}

func NewEnvironment(ID string, requests <-chan string, results chan<- string) *Environment {
	return &Environment{
		ID:       ID,
		Requests: requests,
		Results:  results,
	}
}

func (env *Environment) Run(wg *sync.WaitGroup, stop <-chan bool) {
	fmt.Println("Starting environment " + env.ID)
	for {
		select {
		case request := <-env.Requests:
			fmt.Println(env.ID + " handling " + request)
			if request == "fail" {
				env.Results <- env.ID + " failed [" + request + "]"
				break
			}
			env.Results <- env.ID + " passed [" + request + "]"
		case <-stop:
			wg.Done()
			return
		}
	}
}
