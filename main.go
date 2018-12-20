package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	// Commandline args
	port = flag.String("port", "80", "the port to run the server on")

	// Server config
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"users":["harry","tim","connor"]}`)
	})
	mux.HandleFunc("/tutors", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tutors":["harry"]}`)
	})
	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"students":["tim","connor"]}`)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status":404}`)
	})

	s := &http.Server{
		Addr:         ":" + *port,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Fatal(s.ListenAndServe())
}
