package main

import (
	"flag"
	"github.com/harrisonturton/submission-control/db"
	"github.com/harrisonturton/submission-control/server"
	"log"
	"os"
	"os/signal"
	"sync"
)

var (
	port   = flag.String("port", "80", "the port for the server to run on")
	dbUser = flag.String("dbuser", "harrisonturton", "the database user")
	dbName = flag.String("dbname", "submission_control", "the database name")
)

func main() {
	logger := log.New(os.Stdout, "[server] ", log.LstdFlags)
	// Create store & server
	store, err := db.NewStore(db.Config{
		User:    *dbUser,
		DbName:  *dbName,
		SslMode: "disable",
	})
	failOnError(logger, err)
	srv := server.NewServer(*port, logger, store)
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

func failOnError(logger *log.Logger, err error) {
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
