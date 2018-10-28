// +build !unit

package worker

import (
	"github.com/harrisonturton/submission-control/ci/client"
	"github.com/harrisonturton/submission-control/ci/queue"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	jobQueue    = "job_queue"
	resultQueue = "result_queue"
	testJob     = "hello-world"
	addr        = "amqp://guest:guest@localhost:5672/"
	version     = "1.38"
	timeout     = 5 * time.Second
)

func TestRun(t *testing.T) {
	// Create live worker instance
	client, err := client.New(version)
	if err != nil {
		t.Fatalf("%s", err)
	}
	jobs, err := queue.New(jobQueue, addr)
	failOnError(t, err)
	results, err := queue.New(resultQueue, addr)
	failOnError(t, err)
	worker := New(jobs, results, client, os.Stdout)
	// Send test message
	err = jobs.Push([]byte(testJob))
	failOnError(t, err)
	// Run the worker
	var wg sync.WaitGroup
	done := make(chan bool)
	wg.Add(1)
	go worker.Run(done, &wg)
	go func() {
		defer wg.Done()
		time.Sleep(timeout)
		close(done)
	}()
	wg.Wait()
	// Check that worker popped the message from the job queue
	jobStream := jobs.Stream()
	select {
	case <-jobStream:
		t.Errorf("Failed to remove test message from job queue")
	default:
	}
	// Check the worker put result on the result queue
	resultStream := results.Stream()
	select {
	case <-resultStream:
		break
	default:
		t.Errorf("Failed to add result to result queue")
	}
}

func failOnError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Failed %s", err)
	}
}
