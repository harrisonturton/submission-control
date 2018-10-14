package environment

import (
	"fmt"
	"os"
	"sync"
)

type Environment struct {
	ID       string
	Base     string
	Requests <-chan string
	Results  chan<- string
}

func NewEnvironment(base string, requests <-chan string, results chan<- string) *Environment {
	return &Environment{
		ID:       uuid(),
		Base:     base,
		Requests: requests,
		Results:  results,
	}
}

func (env *Environment) Run(wg *sync.WaitGroup, stop <-chan bool) {
	fmt.Println(fmt.Sprintf("[%s] Starting %s", env.ID, env.Base))
	for {
		select {
		case request := <-env.Requests:
			fmt.Println(fmt.Sprintf("[%s] Handling %s", env.ID, request))
			env.Results <- fmt.Sprintf("[%s] Passed %s", env.ID, request)
		case <-stop:
			wg.Done()
			return
		}
	}
}

func uuid() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x", b[0:4])
}
