package listener

import (
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	resultQueue = "result_queue"
	queueAddr   = "amqp://guest:guest@localhost:5672/"
)

func TestRun(t *testing.T) {
	// Setup listener
	queue := queue.New(5)
	list := New(queue, os.Stdout)
	// Begin testing
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan bool)
	queue.Push([]byte("Test"))
	go list.Run(done, &wg)
	time.Sleep(time.Second * 3)
	close(done)
	wg.Wait()
	// Check results
	select {
	case <-queue.Stream():
		t.Errorf("failed to pop off result queue")
	default:
		return
	}
}
