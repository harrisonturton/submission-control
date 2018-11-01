package queue

import (
	"errors"
)

// Queue is a mocked instance of our RabbitMQ
// connection.
type Queue struct {
	Closed   bool
	Messages chan []byte
}

// New creates a new mocked Queue instance with
// a buffered messages channel of a certain length.
func New(length int) *Queue {
	return &Queue{
		Closed:   false,
		Messages: make(chan []byte, length),
	}
}

// Push will add items to the Messages channel if t
// unless it is full.
func (queue *Queue) Push(data []byte) error {
	select {
	case queue.Messages <- data:
		return nil
	default:
		return errors.New("could not add message, queue buffer is full")
	}
}

// Stream will return a channel where all queue
// items are pushed.
func (queue *Queue) Stream() <-chan []byte {
	return queue.Messages
}

// Length is the number of items in the queue
func (queue *Queue) Length() int {
	return len(Messages)
}

// Close the queue
func (queue *Queue) Close() error {
	queue.Closed = true
	return nil
}
