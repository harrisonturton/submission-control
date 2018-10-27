package queue

import (
	"github.com/streadway/amqp"
	"sync"
)

// Queue is the interface that RabbitMQ must conform to.
// We use this for mocking.
type Queue interface {
	Message(msg string) error
	Consume(wg *sync.WaitGroup, done chan bool, handler func(msg string)) error
}

// RabbitMQ contains the various connections to RabbitMQ
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

// New tries to connect to RabbitMQ on the given address,
// and also creates (if it doesn't already exist) a queue.
func New(name string, addr string) (*RabbitMQ, error) {
	// Attempt to connect
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}
	// Try to open a unique server channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// Declare the queue
	queue, err := ch.QueueDeclare(
		name,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Queue:      &queue,
	}, nil
}

// Message will put a message on the queue. It will not
// close the Connection or Channel, so you must do it
// manually after you've finished messaging.
func (queue *RabbitMQ) Message(msg string) error {
	return queue.Channel.Publish(
		"",               // Exhange
		queue.Queue.Name, // Routing key
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
}

// Consume will continually take mssages from the queue, and give it to the handler function.
// If the done channel is closed during the processing of a message, it will wait until the
// message has finished processing before finishing.
// The queue Connection and Channel will be closed automatically.
func (queue *RabbitMQ) Consume(wg *sync.WaitGroup, done chan bool, handler func(message string)) error {
	defer wg.Done()
	defer queue.Connection.Close()
	defer queue.Channel.Close()
	msgs, err := queue.Channel.Consume(
		queue.Queue.Name, // RabbitMQ
		"",               // Consumer
		true,             // Auto-Ack
		false,            // Exclusive
		false,            // No-local
		false,            // No-Wait
		nil,              // Args
	)
	if err != nil {
		return err
	}
	for {
		select {
		case <-done:
			return nil
		case msg := <-msgs:
			handler(string(msg.Body))
		}
	}
	return nil
}

// Close manually closes the Connection and Channel.
func (queue *RabbitMQ) Close() {
	queue.Connection.Close()
	queue.Channel.Close()
}
