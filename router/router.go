package router

import (
	"fmt"
	"log"
	"net/http"
)

// Router fills the http.ServeHTTP interface, and will
// route individual requests to the nearest handler.
type Router struct {
	logger *log.Logger
	mux    *http.ServeMux
}

// NewRouter creates a new instance of Router attached
// to a Logger instance.
func NewRouter(logger *log.Logger) *Router {
	router := &Router{}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", router.authHandler)
	mux.HandleFunc("/users", router.usersHandler)
	mux.HandleFunc("/", router.notFoundHandler)
	router.mux = mux
	router.logger = logger
	return router
}

// ServeHTTP will route the request to the handler with
// the most similar URL.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Users")
}

func (router *Router) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404")
}
