package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/cache"
	"github.com/harrisonturton/submission-control/ci/producer/listener"
	"github.com/harrisonturton/submission-control/ci/producer/server"
	"github.com/harrisonturton/submission-control/ci/queue"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	serverAddr  = "0.0.0.0:8080"
	queueAddr   = "amqp://guest:guest@rabbitmq:5672/"
	jobQueue    = "job_queue"
	resultQueue = "result_queue"
	pollDelay   = 5 // In seconds
)

func main() {
	// Create server
	jobs := attemptConnect(jobQueue, queueAddr)
	results := attemptConnect(resultQueue, queueAddr)
	cache := cache.New(15, time.Hour*5)
	server := server.New(os.Stdout, jobs, cache, serverAddr)
	listener := listener.New(os.Stdout, results)
	// Run
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(3)
	go server.Serve(done, &wg)
	go listener.Run(done, &wg)
	go func() {
		defer wg.Done()
		<-sig
		fmt.Println("Stopping")
		close(done)
	}()
	wg.Wait()
	fmt.Println("Exiting")
}

func attemptConnect(name string, addr string) queue.Queue {
	for {
		jobs, err := queue.New(name, addr)
		if err != nil {
			fmt.Println("Cannot connect the job queue. Sleeping...")
			time.Sleep(pollDelay * time.Second)
			continue
		}
		return jobs
	}
}
