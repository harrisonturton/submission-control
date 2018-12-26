package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
1) On client, ask for username and password
2) Exchange the username and password for a time-limited access token via HTTPS. Use jwt-go on the server
   to create the token. Use bcrypt to encrypt and compare passwords.
3) Add the recieved access token to the request header for any RESTful API requiring authorization
4) On the server, add an access token checker middleware for those routes. JWT tokens have an expire (exp)
   and not before (nbf) timestamp. JWT validates those when it parses the token from the header.
5) On client, periodically refresh the token. Our tokens expire in 5 minutes. I refresh them every 4 minutes
*/

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

func (router *Router) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		router.notFoundHandler(w, r)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		router.logger.Printf("Error reading request body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	router.logger.Printf("Request: %s\n", string(body))
}

func (router *Router) usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Users")
}

func (router *Router) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404")
}
