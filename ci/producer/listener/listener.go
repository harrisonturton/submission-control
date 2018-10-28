package listener

import (
	"github.com/harrisonturton/submission-control/ci/queue"
	"io"
	"log"
	"sync"
)

// Listener represents a process that continuously listens
// for test results on the result queue, and updates the
// database accordingly.
type Listener struct {
	Results queue.ReadCloser
	Logger  *log.Logger
}

// New creates a new Listener instance.
func New(results queue.ReadCloser, logOut io.Writer) Listener {
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
		case message := <-listener.Results.Stream():
			listener.handleResult(string(message))
		}
	}
}

// handleResult processes every message that the server recieves
// from the result queue.
func (listener *Listener) handleResult(message string) {
	listener.Logger.Printf("Recieved result: %s", message)
}
