package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/producer/server"
	"github.com/harrisonturton/submission-control/ci/queue"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var port = flag.String("port", "8080", "Port for the server to listen on")
var addr = flag.String("addr", "amqp://guest:guest@localhost:5672/", "Address to RabbitMQ service")

const (
	jobQueue    = "job_queue"
	resultQueue = "result_queue"
)

func main() {
	flag.Parse()
	// Create server & queues
	jobs, err := queue.New(jobQueue, *addr)
	panicError(err)
	server := server.New(os.Stdout, jobs, "localhost:"+*port)
	// Listen for SIGINT, CTRL+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	// Run
	var wg sync.WaitGroup
	done := make(chan bool)
	wg.Add(2)
	go server.Serve(done, &wg)
	go func() {
		defer wg.Done()
		<-sig
		fmt.Println("Stopping")
		close(done)
	}()
	wg.Wait()
	fmt.Println("Exiting!")
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
