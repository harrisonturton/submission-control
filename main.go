package main

import (
	"flag"
	"github.com/harrisonturton/submission-control/server"
	"log"
	"os"
	"os/signal"
	"sync"
)

var port = flag.String("port", "80", "port to run the server on")

func main() {
	logger := log.New(os.Stdout, "[server] ", log.LstdFlags)
	logger.Println("Starting...")
	srv := server.NewServer(*port, logger)
	// Boilerplate for graceful shutdown
	var wg sync.WaitGroup
	done := make(chan struct{})
	// Run the server
	wg.Add(1)
	go srv.Serve(&wg, done)
	// Kill server on SIGINT
	wg.Add(1)
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt)
	go func() {
		defer wg.Done()
		<-kill
		close(done)
	}()
	// Wait till killed
	wg.Wait()
	logger.Println("Exiting.")
}
