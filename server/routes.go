package server

import (
	"fmt"
	"net/http"
)

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
