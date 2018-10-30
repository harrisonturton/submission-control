// +build integration
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
	if testing.Short() {
		t.SkipNow()
	}
	// Create a Worker instance
	client, err := client.New(version)
	if err != nil {
		t.Fatalf("%s", err)
	}
	jobs, err := queue.New(jobQueue, addr)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	results, err := queue.New(resultQueue, addr)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	worker := New(jobs, results, client, os.Stdout)
	// Send test message
	err = jobs.Push([]byte(testJob))
	if err != nil {
		t.Fatalf("%s", err)
	}
	// Run the test
	var wg sync.WaitGroup
	done := make(chan bool)
	wg.Add(2)
	go worker.Run(done, &wg)
	go func() {
		defer wg.Done()
		// Wait on result queue, or for timeout
		select {
		case <-results.Stream():
			close(done)
		case <-time.After(timeout):
			close(done)
			t.Fatalf("Failed to add result to result queue")
		}
	}()
	wg.Wait()
	// Check that messages was removed from job queue
	select {
	case <-jobs.Stream():
		t.Fatalf("Failed to remove test message from job queue")
	default:
	}
}
