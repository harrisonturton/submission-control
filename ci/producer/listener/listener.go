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
	Results queue.Queue
	Logger  *log.Logger
}

// New creates a new Listener instance.
func New(results queue.Queue, logOut io.Writer) Listener {
	return Listener{
		Results: results,
		Logger:  log.New(logOut, "", log.LstdFlags),
	}
}

// Run will continuouslly pop messages off the result queue,
// and handle them with listener.handleResult
func (listener *Listener) Run(done chan bool, wg *sync.WaitGroup) {
	listener.Logger.Printf("Waiting for results on queue %s", listener.Results.Name())
	listener.Results.Consume(wg, done, listener.handleResult)
}

// handleResult processes every message that the server recieves
// from the result queue.
func (listener *Listener) handleResult(message string) {
	listener.Logger.Printf("Recieved result: %s", message)
}
