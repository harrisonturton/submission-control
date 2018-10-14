package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/cli"
	"github.com/harrisonturton/submission-control/environment"
	"sync"
)

const (
	EnvironmentCount = 5
	MessageCount     = 10
)

var wg sync.WaitGroup
var requests = make(chan string, MessageCount)
var results = make(chan string, MessageCount)
var stop = make(chan bool)

func main() {
	for i := 0; i < EnvironmentCount; i++ {
		wg.Add(1)
		env := environment.NewEnvironment("python-base", requests, results)
		go env.Run(&wg, stop)
	}
	wg.Add(2)
	go cli.Run(stop, &wg, requests)
	go func() {
		for {
			select {
			case result := <-results:
				fmt.Println(result)
			case <-stop:
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	fmt.Println("Goodbye!")
}
