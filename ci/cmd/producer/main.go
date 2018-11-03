package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/cache"
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
)

func main() {
	server := createServer()
	// Get notified upon SIGINT
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	// Run server
	var wg sync.WaitGroup
	done := make(chan bool)
	wg.Add(2)
	go server.Serve(done, &wg)
	go func() {
		defer wg.Done()
		<-sig
		fmt.Println("Stopping")
		close(done)
	}()
	wg.Wait()
	fmt.Println("Exiting")
}

// Will continuously try and connect until success
func createServer() *server.Server {
	for {
		jobs, err := queue.New(jobQueue, queueAddr)
		if err != nil {
			fmt.Println("%s\nTrying again in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		cache := cache.New(15, time.Hour*5)
		return server.New(os.Stdout, jobs, cache, serverAddr)
	}
}

func exitError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
