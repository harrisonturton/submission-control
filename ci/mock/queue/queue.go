package queue

import (
	"errors"
	"sync"
)

// Queue is a mocked instance of our RabbitMQ
// connection.
type Queue struct {
	QueueName string
	Messages  chan string
}

// New creates a new mocked Queue instance.
func New(name string, addr string) (*Queue, error) {
	return &Queue{name, make(chan string, 10)}, nil
}

// Name will return the name of the queue
func (queue *Queue) Name() string {
	return queue.QueueName
}

// Message will send a dummy message to the queue. If
// the channel is full, then it will return an error.
func (queue *Queue) Message(msg string) error {
	select {
	case queue.Messages <- msg:
		return nil
	default:
		return errors.New("could not add message")
	}
}

// Consume will handle everything on the channel through the handle func,
// and will only stop when the done channel recieves input.
func (queue *Queue) Consume(wg *sync.WaitGroup, done chan bool, handler func(msg string)) error {
	defer wg.Done()
	for {
		select {
		case <-done:
			return nil
		case msg := <-queue.Messages:
			handler(string(msg))
		}
	}
}
