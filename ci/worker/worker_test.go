// +build unit
package worker

import (
	"github.com/harrisonturton/submission-control/ci/mock/client"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"testing"
)

func TestHandleJob(t *testing.T) {
	t.Parallel()
	// Create worker
	client := client.New()
	jobs := queue.New(5)
	results := queue.New(5)
	worker := New(jobs, results, client, os.Stdout)
	// Test handleJob
	worker.handleJob("test")
	// Make sure something is put on result queue
	select {
	case <-results.Messages:
		break
	default:
		t.Fatalf("Failed to put data on result queue")
	}
	if results.Closed || jobs.Closed {
		t.Fatalf("Prematurely closed the queues")
	}
}
