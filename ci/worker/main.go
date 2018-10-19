package main

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/worker/worker"
	"os"
	"sync"
	"time"
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
	done := make(chan bool)
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 10)
		close(done)
	}()
	go worker.Run(done, &wg)
	wg.Wait()

	fmt.Println("Finishing.")
}
