package server

import (
	"context"
	"fmt"
	"github.com/harrisonturton/submission-control/test-engine/queue"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Server   *http.Server
	Logger   *log.Logger
	JobQueue *queue.Queue
}

const (
	JobQueue = "job_queue"
)

// New creates a new Server instance. It will return and error
// if it cannot connect to the RabbitMQ queues.
func New(logOut io.Writer, addr string) (*Server, error) {
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
	}
	router := http.NewServeMux()
	router.HandleFunc("/", server.handleRequest)
	server.Server.Handler = router
	jobQueue, err := queue.New(JobQueue, "amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	server.JobQueue = jobQueue
	return server, nil
}

// Serve will listen for requests on the Server address. It will exit
// with exit code 1 if it encouters an error.
func (server *Server) Serve(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		<-done
		server.Logger.Println("Attempting to shutdown server...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

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
