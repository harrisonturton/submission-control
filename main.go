package main

import (
	"flag"
	api "github.com/harrisonturton/submission-control/server"
	"log"
	"os"
	"os/signal"
	"sync"
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
	wg.Add(1)
	go server.Serve(logger, &wg, done)
	// Kill server on SIGINT
	wg.Add(1)
	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt)
	go func() {
		defer wg.Done()
		<-killChan
		close(done)
	}()
	// Wait till killed
	wg.Wait()
	logger.Println("Exiting.")
}
