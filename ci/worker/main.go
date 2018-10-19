package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/worker/worker"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	JobQueue    = "job_queue"
	ResultQueue = "result_queue"
)

func main() {
	worker, err := worker.New(os.Stdout, JobQueue, ResultQueue, "amqp://guest:guest@localhost:5672/")
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
}
