package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/producer/listener"
	"github.com/harrisonturton/submission-control/ci/producer/server"
	"github.com/harrisonturton/submission-control/ci/queue"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var port = flag.String("port", "8080", "Port for the server to listen on")
var addr = flag.String("addr", "amqp://guest:guest@rabbitmq:5672/", "Address to RabbitMQ service")

const (
	jobQueue    = "job_queue"
	resultQueue = "result_queue"
)

func main() {
	flag.Parse()
	// Create server & queues
	jobs, err := queue.New(jobQueue, *addr)
	exitError(err)
	results, err := queue.New(resultQueue, *addr)
	exitError(err)
	server := server.New(os.Stdout, jobs, "localhost:"+*port)
	list := listener.New(results, os.Stdout)
	// Listen for SIGINT, CTRL+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	// Run
	var wg sync.WaitGroup
	done := make(chan bool)
	wg.Add(3)
	go server.Serve(done, &wg)
	go list.Run(done, &wg)
	go func() {
		defer wg.Done()
		<-sig
		fmt.Println("Stopping")
		close(done)
	}()
	wg.Wait()
	fmt.Println("Exiting!")
}

func exitError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
