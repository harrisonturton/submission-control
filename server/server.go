package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server represents the HTTP server, with added
// graceful shutdown and middleware & tracing functionality.
type Server struct {
	logger     *log.Logger
	httpServer *http.Server
	middleware []func(http.Handler) http.Handler
}

var (
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 15 * time.Second
)

// NewServer creates a new Server instance,
// and attaches all required tracing.
func NewServer(port string, logger *log.Logger) *Server {
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      traceRoutes(logger),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

// Serve will listen for requests on the given port. It will output
// any errors to the logger.
func (server *Server) Serve(wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	go server.waitForShutdown(done)
	server.logger.Printf("Server starting on %s\n", server.httpServer.Addr)
	err := server.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		server.logger.Fatalf("Server failed: %s\n", err.Error())
	}
	server.logger.Println("Server stopped.")
}

// waitForShutdown will wait until the done channel is closed,
// and attempt to gracefully shutdown the server. If it doesn't
// shutdown within the shutdownTimeout, it will forcefully shutdown.
// Graceful shutdown waits for all outstanding requests to be handled.
func (server *Server) waitForShutdown(done chan struct{}) {
	<-done
	server.logger.Println("Attempting to shutdown...")
	// Forcefully shutdown after shutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	// Attempt to gracefully shutdown
	err := server.httpServer.Shutdown(ctx)
	if err != nil {
		server.logger.Fatalf("Could not gracefully shutdown: %s\n", err.Error())
	}
}

// traceRoutes builds the routes (a http.Handler) with
// some tracing to log important info about each request.
func traceRoutes(logger *log.Logger) http.Handler {
	router := NewRouter(logger)
	return trace(logger, router)
}

// trace will intercept each request, log some information about it,
// and then pass it along to the proper handlers.
func trace(logger *log.Logger, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		printRequest(logger, r)
		nextHandler.ServeHTTP(w, r)
	})
}

func printRequest(logger *log.Logger, r *http.Request) {
	if r.URL != nil {
		logger.Printf("[%s] %s\n", r.Method, r.URL)
	} else {
		logger.Printf("[%s] [no url]\n", r.Method)
	}
	logger.Printf("Content Length: %d\n", r.ContentLength)
}
