package server

import (
	"fmt"
	"log"
	"net/http"
)

var routes = map[string]http.HandlerFunc{
	"/users":    usersHandler,
	"/tutors":   tutorsHandler,
	"/students": studentsHandler,
	"/":         notFoundHandler,
}

// makeRoutes creates a http.ServeMux from the routes constant
func makeRoutes(logger *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	for route, handler := range routes {
		mux.HandleFunc(route, handler)
	}
	return mux
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Users")
}

func tutorsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tutors")
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Students")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404")
}
