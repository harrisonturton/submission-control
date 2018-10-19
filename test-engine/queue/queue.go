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
	defer conn.Close()
	// Try to open a unique server channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()
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

// Message will put a message on the queue
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

func (queue *Queue) Consume(wg *sync.WaitGroup, done chan bool, handler func(message string)) error {
	defer wg.Done()
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

	for msg := range msgs {
		select {
		case <-done:
			return nil
		default:
			handler(string(msg.Body))
		}
	}
	return nil
}
