package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/submission-control/test-engine/producer/listener"
	"github.com/harrisonturton/submission-control/test-engine/producer/server"
	"os"
	"sync"
	"time"
)

var port = flag.String("port", "8080", "Port for the server to listen on.")
var addr = flag.String("addr", "amqp://guest:guest@localhost:5672/", "Address to RabbitMQ service")

const ResultQueue = "result_queue"

func main() {
	flag.Parse()
	server, err := server.New(os.Stdout, "localhost:"+*port)
	panicError(err)
	listener, err := listener.New(os.Stdout, ResultQueue, *addr)
	panicError(err)

	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(3)
	go server.Serve(done, &wg)
	go listener.Run(done, &wg)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 100)
		fmt.Println("Stopping...")
		close(done)
	}()
	wg.Wait()

	fmt.Println("Exiting")
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
