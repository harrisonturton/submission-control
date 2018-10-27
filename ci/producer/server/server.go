package server

import (
	"context"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/queue"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server represents the goroutine that listens
// for test requests over the web, and puts them
// on the job queue.
type Server struct {
	Server *http.Server
	Logger *log.Logger
	Jobs   queue.Queue
}

const (
	jobQueue = "job_queue"
)

// New creates a new Server instance.
func New(logOut io.Writer, jobs queue.Queue, addr string) *Server {
	logger := log.New(logOut, "", log.LstdFlags)
	server := &Server{
		Server: &http.Server{
			Addr:         addr,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		Logger: logger,
		Jobs:   jobs,
	}
	router := http.NewServeMux()
	router.HandleFunc("/", server.handleRequest)
	server.Server.Handler = withLogging(logger, router)
	return server
}

// Serve will continuously listen for requests, and handle them
// with server.handleRequest.
func (server *Server) Serve(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		<-done
		// Force shutdown if couldn't gracefully shutdown within the timeout
		server.Logger.Println("Attempting to shutdown server...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		// Attempt to gracefully shutdown
		err := server.Server.Shutdown(ctx)
		if err != nil {
			server.Logger.Fatalf("Could not gracefully shutdown the server. %s\n", err.Error())
		}
	}()
	server.Logger.Printf("Server starting at %s\n", server.Server.Addr)
	if err := server.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		server.Logger.Fatalf("Server failed. %s", err.Error())
	}
	server.Logger.Printf("Server stopped.")
}

// handleRequest handles every request to come through the server
func (server *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	server.Logger.Printf("%s: %s", r.Method, r.URL.Path)
	server.Logger.Printf("From: %s", r.RemoteAddr)
	server.Logger.Printf("As: %s", r.UserAgent())
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error processing request body.")
	} else {
		fmt.Fprintf(w, "Handled!")
	}
}

func withLogging(logger *log.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s: %s", r.Method, r.URL.Path)
		logger.Printf("From: %s", r.RemoteAddr)
		logger.Printf("As: %s", r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}
