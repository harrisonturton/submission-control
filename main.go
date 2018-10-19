package main

import (
	"flag"
	"github.com/streadway/amqp"
	"log"
)

const (
	QueueName = "job_queue"
)

var queueAddr = flag.String("addr", "amqp://guest:guest@localhost:5672/", "Job queue address")

func main() {
	flag.Parse()

	log.Println("Attempting to connect to " + *queueAddr)
	conn, err := amqp.Dial(*queueAddr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel.")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		QueueName, // Name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	failOnError(err, "Failed to declare queue: "+QueueName)

	msgs, err := ch.Consume(
		queue.Name, // Queue
		"",         // Consumer
		true,       // Auto-Ack
		false,      // Exclusive
		false,      // No-local
		false,      // No-Wait
		nil,        // Args
	)
	failOnError(err, "Failed to register a consumer")

	// Handle messages forever
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Recieved message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// failOnError will print the error & message before
// exiting with exit code 1.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
