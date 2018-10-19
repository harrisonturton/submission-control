package queue

import (
	"github.com/streadway/amqp"
	"sync"
)

type Queue struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

// New tries to connect to RabbitMQ on the given address,
// and also creates (if it doesn't already exist) a queue.
func New(name string, addr string) (*Queue, error) {
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
	return &Queue{
		Connection: conn,
		Channel:    ch,
		Queue:      &queue,
	}, nil
}

// Message will put a message on the queue. It will not
// close the Connection or Channel, so you must do it
// manually after you've finished messaging.
func (queue *Queue) Message(msg string) error {
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
// The queue Connection and Channel will be closed after this function.
func (queue *Queue) Consume(wg *sync.WaitGroup, done chan bool, handler func(message string)) error {
	defer wg.Done()
	defer queue.Connection.Close()
	defer queue.Channel.Close()
	msgs, err := queue.Channel.Consume(
		queue.Queue.Name, // Queue
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
func (queue *Queue) Close() {
	queue.Connection.Close()
	queue.Channel.Close()
}
