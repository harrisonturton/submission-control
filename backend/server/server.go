package server

import (
	"context"
	"github.com/harrisonturton/submission-control/routes"
	"github.com/harrisonturton/submission-control/store"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server listens and responds to API requests. It wraps a
// http.Server instance to provide graceful shutdown, logging
// and middleware functionality.
type Server struct {
	logger *log.Logger
	server *http.Server
	store  *store.Store
}

var (
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 15 * time.Second
)

// NewServer creates a new Server instance.
func NewServer(port string, logger *log.Logger, store *store.Store) *Server {
	mux := routes.CreateMux(store)
	handler := addMiddleware(
		mux,
		logAll(logger),
		attachContext(),
	)
	return &Server{
		logger: logger,
		store:  store,
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

// Serve will start running the server, starting to listen for requests on
// the given port. It will gracefully shutdown when the done channel is closed.
func (server *Server) Serve(wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	go server.waitForShutdown(done)
	server.logger.Printf("Server starting on %s\n", server.server.Addr)
	err := server.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		server.logger.Fatalf("Server failed: %v\n", err)
	}
	server.logger.Println("Server stopped.")
}

// waitForShutdown will wait until the done channel is closed,
// and attempt to gracefully shutdown the server. It will try to
// serve all remaining requests before stopping.
// If it cannot finish serving the requests within the shutdownTimeout,
// it will forcefully stop.
func (server *Server) waitForShutdown(done chan struct{}) {
	<-done
	server.logger.Println("Attempting to gracefully shutdown...")
	// Forcefully shutdown after shutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	// Attempt to gracefully shutdown
	if err := server.server.Shutdown(ctx); err != nil {
		server.logger.Fatalf("Could not gracefully shutdown: %v\n", err)
		return
	}
	server.logger.Println("Stopped gracefully!")
}
