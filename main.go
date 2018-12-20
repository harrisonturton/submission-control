package main

import (
	"flag"
	api "github.com/harrisonturton/submission-control/server"
	"log"
	"os"
	"sync"
	"time"
)

// Commandline args
var port = flag.String("port", "80", "the port to run the server on")

func main() {
	logger := log.New(os.Stdout, "[server] ", log.LstdFlags)
	logger.Println("Starting...")
	server := api.NewServer(*port)
	// Boilerplate for graceful shutdown
	var wg sync.WaitGroup
	done := make(chan struct{})
	// Run the server
	wg.Add(2)
	go server.Serve(logger, &wg, done)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Second)
		close(done)
	}()
	wg.Wait()
	logger.Println("Exiting.")
}
