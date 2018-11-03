package listener

import (
	"github.com/harrisonturton/submission-control/ci/queue"
	"github.com/harrisonturton/submission-control/ci/types"
	"io"
	"log"
	"sync"
)

// Listener represents a process that continuously listens
// for test results on the result queue, and updates the
// database accordingly.
type Listener struct {
	Results queue.Queue
	Logger  *log.Logger
}

// New creates a new Listener instance.
func New(logOut io.Writer, results queue.Queue) Listener {
	return Listener{
		Results: results,
		Logger:  log.New(logOut, "", log.LstdFlags),
	}
}

// Run will continuouslly pop messages off the result queue,
// and handle them with listener.handleResult
func (listener *Listener) Run(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	listener.Logger.Printf("Waiting for results ...")
	for {
		select {
		case <-done:
			return
		case resultGob := <-listener.Results.Stream():
			result := types.TestResult{}
			err := result.Deserialize(resultGob)
			if err != nil {
				listener.Logger.Println("Failed to deserialize TestResult gob.")
			}
			listener.handleResult(result)
		}
	}
}

// handleResult processes every message that the server recieves
// from the result queue.
func (listener *Listener) handleResult(result types.TestResult) {
	listener.Logger.Printf("Recieved result with stdout: %s\nstderr: %s", result.Stdout, result.Stderr)
}
