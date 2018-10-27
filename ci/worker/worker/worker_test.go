package worker

import (
	"github.com/harrisonturton/submission-control/ci/mock/client"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"sync"
	"testing"
)

const (
	jobQueueName    = "job_queue"
	resultQueueName = "result_queue"
	queueAddr       = "amqp://guest:guest@localhost:5672/"
	dockerVersion   = "1.38"
)

func TestNew(t *testing.T) {
	client, _ := client.New(dockerVersion)
	jobQueue, _ := queue.New(jobQueueName, queueAddr)
	resultQueue, _ := queue.New(resultQueueName, queueAddr)
	New(jobQueue, resultQueue, client, os.Stdout)
}

func TestRun(t *testing.T) {
	// Create worker
	client, _ := client.New(dockerVersion)
	jobs, _ := queue.New(jobQueueName, queueAddr)
	results, _ := queue.New(resultQueueName, queueAddr)
	worker := New(jobs, results, client, os.Stdout)
	// Test
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan bool)
	jobs.Message("test")
	go worker.Run(done, &wg)
	close(done)
	wg.Wait()
	select {
	case <-results.Messages:
		break
	default:
		t.Errorf("Failed to put data on result queue")
	}
}

func TestHandleJob(t *testing.T) {
	// Create worker
	client, _ := client.New(dockerVersion)
	jobs, _ := queue.New(jobQueueName, queueAddr)
	results, _ := queue.New(resultQueueName, queueAddr)
	worker := New(jobs, results, client, os.Stdout)
	// Test handleJob
	worker.handleJob("test")
	// Make sure correct stuff is put on result queue
	select {
	case <-results.Messages:
		break
	default:
		t.Errorf("Failed to put data on result queue")
	}
}
