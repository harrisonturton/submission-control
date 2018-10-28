package worker

import (
	"github.com/harrisonturton/submission-control/ci/mock/client"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"testing"
)

const (
	jobQueue      = "job_queue"
	resultQueue   = "result_queue"
	queueAddr     = "amqp://guest:guest@localhost:5672/"
	dockerVersion = "1.38"
)

func TestHandleJob(t *testing.T) {
	// Create worker
	client, _ := client.New(dockerVersion)
	jobs := queue.New(5)
	results := queue.New(5)
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
	if results.Closed || jobs.Closed {
		t.Errorf("Prematurely closed the queues")
	}
}
