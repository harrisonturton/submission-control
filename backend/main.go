package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"database/sql"
	"github.com/harrisonturton/submission-control/backend/store"
	_ "github.com/lib/pq"

	"github.com/harrisonturton/submission-control/backend/server"
	"os/signal"
	"sync"
)

var (
	port   = flag.String("port", "80", "the port for the server to run on")
	dbUser = flag.String("dbuser", "harrisonturton", "the database user")
	dbName = flag.String("dbname", "submission_control", "the database name")
)

const (
	certFile = "server.crt"
	keyFile  = "server.key"
)

func main() {
	logger := log.New(os.Stdout, "[server] ", log.LstdFlags)
	store := createStore(logger)
	srv := server.NewServer(*port, logger, store)
	runServer(srv)
	logger.Println("Exiting.")
}

// createStore will create a database connection and
// a store around that.
func createStore(logger *log.Logger) *store.Store {
	db, err := createDatabase()
	if err != nil {
		logger.Fatal(err)
	}
	store, err := store.NewStore(db)
	if err != nil {
		logger.Fatal(err)
	}
	return store
}

// createDatabase will try and open a new connection to the
// database (but potentially re-uses an old one, due to
// connection pooling)
func createDatabase() (*sql.DB, error) {
	conn := fmt.Sprintf("user=%s dbname=%s sslmode=disable", *dbUser, *dbName)
	return sql.Open("postgres", conn)
}

// runServer will run the server until an interrupt is sent
func runServer(srv *server.Server) {
	// Boilerplate for graceful shutdown
	var wg sync.WaitGroup
	done := make(chan struct{})
	// Run server
	wg.Add(1)
	go srv.ServeTLS(certFile, keyFile, &wg, done)
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
}
