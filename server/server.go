package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Instance http.Server
}

var timeout = time.Second * 10

// Create a new Server instance
func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/exit", handleExit)
	return &Server{http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
	}}
}

// Run the server until the stop channel is closed.
func (server *Server) Serve(stop chan bool, wg *sync.WaitGroup) {
	fmt.Println("Server starting on " + server.Instance.Addr)
	defer wg.Done()
	go func() {
		<-stop
		fmt.Println("Shutting down server...")
		if err := server.Instance.Shutdown(context.Background()); err != nil {
			// Error on closing listeners, or context timeout
			fmt.Println("Error shutting down: " + err.Error())
		}
	}()
	if err := server.Instance.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("Error listening and serving: " + err.Error())
	}
	fmt.Println("Server shut down.")
}

// Handle requests made to the root directory
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello root!")
}

// Handle requests made to the /exit directory
func handleExit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello exit!")
}
