package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/submission-control/test-engine/producer/server"
	"os"
	"sync"
	"time"
)

var port = flag.String("port", "8080", "Port for the server to listen on.")

func main() {
	flag.Parse()
	server, err := server.New(os.Stdout, "localhost:"+*port)
	panicError(err)

	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go server.Serve(done, &wg)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 10)
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
