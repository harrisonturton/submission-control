package server

import (
	"context"
	"fmt"
	"github.com/harrisonturton/submission-control/worker/client"
	"github.com/harrisonturton/submission-control/worker/environment"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Server       *http.Server
	Client       *client.Client
	Logger       *log.Logger
	Environments map[string]*environment.Environment
}

// NewServer creates a new load server instance. If the Docker worker API version
// is incompatible with the Go SDK, try passing in a different version number.
func New(version string, addr string, images []string, logOut io.Writer) (Server, error) {
	client, err := client.NewClient(version)
	if err != nil {
		return Server{}, err
	}
	logger := log.New(logOut, "", log.LstdFlags)
	envs := make(map[string]*environment.Environment)
	for _, image := range images {
		envs[image] = environment.New(image, []string{}, client, logger)
	}
	server := Server{
		Server: &http.Server{
			Addr:         addr,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		Client:       client,
		Logger:       logger,
		Environments: envs,
	}
	router := http.NewServeMux()
	router.HandleFunc("/", server.handleRequest)
	server.Server.Handler = middlewareLogs(router, logger)
	return server, nil
}

// Serve will listen for requests on the Server address. Warning: this
// will panic if it encounters an error.
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

// handlRequest handles every incoming request. It then passes the data
// to the relevant environment.
func (server *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handled!")
}

// middlewareLogs is middleware that lies between when a request is recieved,
// and when it is handled. It logs basic info about the request.
func middlewareLogs(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s: %s", r.Method, r.URL.Path)
			logger.Printf("From: %s\nAs: %s", r.RemoteAddr, r.UserAgent())
			next.ServeHTTP(w, r)
		})
}
