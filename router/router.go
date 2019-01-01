package router

import (
	"fmt"
	"github.com/harrisonturton/submission-control/request"
	"github.com/harrisonturton/submission-control/store"
	"log"
	"net/http"
)

// Router fills the http.ServeHTTP interface, and will
// route individual requests to the nearest handler.
type Router struct {
	logger *log.Logger
	store  *store.Store
	mux    *http.ServeMux
}

// NewRouter creates a new instance of Router attached
// to a Logger instance.
func NewRouter(logger *log.Logger, store *store.Store) *Router {
	router := &Router{}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", router.authHandler)
	mux.HandleFunc("/users", router.usersHandler)
	mux.HandleFunc("/", router.notFoundHandler)
	router.mux = mux
	router.logger = logger
	router.store = store
	return router
}

// ServeHTTP will route the request to the handler with
// the most similar URL.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) usersHandler(w http.ResponseWriter, r *http.Request) {
	if !request.GetAuth(r) {
		fmt.Fprintf(w, "authentication error")
		router.logger.Println("Could not authenticate")
		return
	}
	fmt.Fprintf(w, "Users")
}

func (router *Router) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404")
}
