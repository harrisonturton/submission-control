package main

import (
<<<<<<< HEAD
	"bytes"
	"flag"
	"github.com/harrisonturton/submission-control/worker/client"
	"github.com/streadway/amqp"
	"log"
	"time"
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

	client, err := client.New("1.38")
	failOnError(err, "Failed to create Docker client.")

	// Handle messages forever
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			handleMessage(client, string(d.Body))
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func handleMessage(client *client.Client, msg string) {
	log.Printf("Recieved message: %s", msg)
	resp, err := client.CreateContainer(msg, []string{})
	if err != nil {
		log.Printf("Failed to create container: %s", err)
		return
	}

	client.WaitForContainer(resp.ID, time.Second*30)
	logReader, err := client.ReadContainerLogs(resp.ID, true, true)
	buf := new(bytes.Buffer)
	buf.ReadFrom(logReader)
	log.Printf("Result: %s", buf.String())
}

// failOnError will print the error & message before
// exiting with exit code 1.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
=======
	"fmt"
	"github.com/harrisonturton/hydra-cli/handlers"
	"net/rpc"
	"os"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":3000")
	panicError(err)
	err = handlers.RunCommand(os.Args, client)
	panicError(err)
}

func panicError(err error) {
	if err != nil {
		fmt.Println(err.Error())
>>>>>>> 53f405fc1f5594fd337c0f518392b6655a7c851a
	}
}
