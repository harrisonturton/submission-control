package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/environment"
	"strconv"
	"sync"
)

const (
	EnvironmentCount = 5
	MessageCount     = 10
	Delay            = 500
)

var requests = make(chan string, MessageCount)
var results = make(chan string, MessageCount)
var stop = make(chan bool, EnvironmentCount)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < EnvironmentCount; i++ {
		wg.Add(1)
		env := environment.NewEnvironment(strconv.Itoa(i), requests, results)
		go env.Run(&wg, stop)
	}
	go func() {
		for i := 0; i < MessageCount; i++ {
			requests <- "request " + strconv.Itoa(i)
		}
		for i := 0; i < MessageCount; i++ {
			result := <-results
			fmt.Println(result)
		}
		for i := 0; i < EnvironmentCount; i++ {
			stop <- true
		}
	}()
	wg.Wait()
	fmt.Println("Done!")
}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
