package listener

import (
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"sync"
	"testing"
)

const (
	resultQueue = "result_queue"
	queueAddr   = "amqp://guest:guest@localhost:5672/"
)

func TestRun(t *testing.T) {
	// Setup listener
	queue, _ := queue.New(resultQueue, queueAddr)
	list := New(queue, os.Stdout)
	// Begin testing
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan bool)
	queue.Message("Test")
	go list.Run(done, &wg)
	close(done)
	wg.Wait()
	// Check results
	select {
	case <-queue.Messages:
		t.Errorf("failed to pop off result queue")
	default:
		return
	}
}
