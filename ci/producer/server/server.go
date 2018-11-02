package server

import (
	"context"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/cache"
	"github.com/harrisonturton/submission-control/ci/queue"
	"github.com/harrisonturton/submission-control/ci/types"
	"io"
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
	Cache  *cache.Cache
}

const shutdownTimeout = 15 * time.Second

// New creates a new Server instance.
func New(logOut io.Writer, jobs queue.Queue, cache *cache.Cache, addr string) *Server {
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
		Cache:  cache,
	}
	router := http.NewServeMux()
	router.HandleFunc("/", server.handleRequest)
	server.Server.Handler = withLogging(logger, router)
	return server
}

// Serve will continuously serve requests, either responding with
// test results or putting new jobs on the queue.
func (server *Server) Serve(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		<-done
		// Force shutdown if couldn't gracefully shutdown before the timeout
		server.Logger.Println("Attempting to shutdown server...")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
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
	if r.Method != http.MethodPost {
		return
	}
	defer r.Body.Close()
	config := types.TestConfig{}
	err := config.UnmarshalJSON(r.Body)
	if err != nil {
		server.Logger.Printf("Failed to unmarshal request body: %s\n", err)
		w.Write([]byte(fmt.Sprintf("Failed to unmarshal request: %s\n", err)))
		return
	}
	server.Logger.Printf("Got job with config version %s and image %s\n", *config.Version, *config.Env.Image)
	bytes, err := config.Serialize()
	if err != nil {
		server.Logger.Printf("Failed to serialize job")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to serialize request: %s\n", err)))
		return
	}
	err = server.Jobs.Push(bytes)
	if err != nil {
		server.Logger.Printf("Failed to push job to job queue.")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error."))
		return
	}
	w.Write([]byte("Got the job!\n"))
}

// withLogging is middleware that wraps a http Handler. It logs basic info
// about every request that passes through it.
func withLogging(logger *log.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s: %s", r.Method, r.URL.Path)
		logger.Printf("From: %s", r.RemoteAddr)
		logger.Printf("As: %s", r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}
