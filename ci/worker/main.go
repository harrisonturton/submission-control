package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/client"
	"github.com/harrisonturton/submission-control/ci/queue"
	"github.com/harrisonturton/submission-control/ci/worker/worker"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	jobQueueName    = "job_queue"
	resultQueueName = "result_queue"
	host            = "rabbitmq"
	queueAddr       = "amqp://guest:guest@" + host + ":5672/"
	dockerVersion   = "1.38"
)

func main() {
	client, err := client.New(dockerVersion)
	exitError(err)
	jobQueue, err := queue.New(jobQueueName, queueAddr)
	exitError(err)
	resultQueue, err := queue.New(resultQueueName, queueAddr)
	exitError(err)
	worker := worker.New(jobQueue, resultQueue, client, os.Stdout)

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

func exitError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

/*
func main() {
	worker, err := worker.New(os.Stdout, jobQueue, resultQueue, "amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

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

	fmt.Println("Finished.")
}*/
