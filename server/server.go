package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server handles all API requests and passes
// the data to the correct handlers.
type Server struct {
	Server *http.Server
}

var (
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 15 * time.Second
)

// NewServer creates a new instance of Server, but
// does not begin running it.
func NewServer(port string) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         ":" + port,
			Handler:      makeMux(),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

// Serve will run the server and output any errors
// to the logger.
func (server *Server) Serve(logger *log.Logger, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	go func() {
		<-done
		// Force shutdown if couldn't gracefully shutdown
		logger.Println("Attempting to shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		// Attempt to gracefully shutdown
		if err := server.Server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown: %s\n", err.Error())
		}
	}()
	logger.Printf("Server starting on %s\n", server.Server.Addr)
	err := server.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed: %s\n", err.Error())
	}
	logger.Println("Server stopped.")
}

func makeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/tutors", tutorsHandler)
	mux.HandleFunc("/students", studentsHandler)
	mux.HandleFunc("/", notFoundHandler)
	return mux
}
