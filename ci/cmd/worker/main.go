package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/client"
	"github.com/harrisonturton/submission-control/ci/queue"
	"github.com/harrisonturton/submission-control/ci/worker"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	jobQueueName    = "job_queue"
	resultQueueName = "result_queue"
	queueAddr       = "amqp://guest:guest@rabbitmq:5672/"
	dockerVersion   = "1.38"
)

func main() {
	worker := createWorker()
	fmt.Println("WORKER RUNNING")

	var wg sync.WaitGroup
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	done := make(chan bool)
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-sig
		fmt.Println("Finishing...")
		close(done)
	}()
	go worker.Run(done, &wg)
	wg.Wait()
	fmt.Println("Finished")
}

func createWorker() *worker.Worker {
	for {
		client, err := client.New(dockerVersion)
		if err != nil {
			fmt.Println("Worker sleeping for 5 seconds: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		jobQueue, err := queue.New(jobQueueName, queueAddr)
		if err != nil {
			fmt.Println("Worker sleeping for 5 seconds: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		resultQueue, err := queue.New(resultQueueName, queueAddr)
		if err != nil {
			fmt.Println("Worker sleeping for 5 seconds: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return worker.New(jobQueue, resultQueue, client, os.Stdout)
	}
}

func exitError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
