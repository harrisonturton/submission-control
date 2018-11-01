package queue

import (
	"github.com/streadway/amqp"
	"sync"
)

// Queue represents a connection to RabbitMQ.
type Queue interface {
	Push(data []byte) error
	Stream() <-chan []byte
	Close() error
	Length() int
}

// RabbitMQ contains the various connections to RabbitMQ
type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	messages   chan []byte
	wg         sync.WaitGroup
	done       chan bool
}

// New will create a new RabbitMQ queue instance. It will
// declare the queue, which creates a new instance if one
// doesn't already exist.
func New(name string, addr string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
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
	result := &RabbitMQ{
		connection: conn,
		channel:    ch,
		queue:      &queue,
		messages:   make(chan []byte),
		done:       make(chan bool),
	}
	result.wg.Add(1)
	go result.listen()
	return result, nil
}

// Push will put an item on the queue.
func (queue *RabbitMQ) Push(data []byte) error {
	return queue.channel.Publish(
		"",               // Exhange
		queue.queue.Name, // Routing key
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

// Stream will continuously put queue items on the channel.
func (queue *RabbitMQ) Stream() <-chan []byte {
	return queue.messages
}

// Close closes the channel and connection to RabbitMQ.
// May block indefinitely if queue.Listen never finishes.
func (queue *RabbitMQ) Close() error {
	close(queue.done)
	queue.wg.Wait()
	return queue.connection.Close()
}

// Length is the number of items in the queue.
func (queue *RabbitMQ) Length() int {
	return queue.queue.Messages
}

// Listen will convert incoming queue messages into []byte
// items, and put them on the messages channel if there is a
// consumer.
func (queue *RabbitMQ) listen() error {
	defer queue.wg.Done()
	messages, err := queue.channel.Consume(
		queue.queue.Name,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
	if err != nil {
		return err
	}
	for {
		select {
		case <-queue.done:
			return nil
		case data := <-messages:
			queue.messages <- data.Body
			data.Ack(false)
		}
	}
}
