package routes

import (
	"encoding/json"
	"fmt"
	"github.com/harrisonturton/submission-control/auth"
	"github.com/harrisonturton/submission-control/request"
	"github.com/harrisonturton/submission-control/store"
	"net/http"
)

// These are the handlers called for each route, as specified in
// routes/router.go
//
// http.ServeMux expects a http.Handler, which is an interface for types
// with a ServeHTTP(http.ResponseWriter, *http.Request) function.
// Since this does not give any space to add extra dependencies (e.g.
// logger and database instances), I've manipulated closures to give the
// handler access to these.
//
// I could have implemented a custom Router{} type, with fields for various
// database and logger instances, but this is slightly cleaner and (hopefully)
// easier to test.

// authHandler is called on the /auth route to request a new JWT token.
// It will authenticate the LoginRequest and generate a new token.
func authHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reject if not POST
		if r.Method != http.MethodPost {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Unmarshal POST data
		var login LoginRequest
		err := json.Unmarshal(request.GetBody(r), &login)
		if err != nil {
			http.Error(w, "unrecognised body", http.StatusBadRequest)
			return
		}
		// Verify login information
		ok, err := auth.Authenticate(store, login.Email, login.Password)
		if !ok || err != nil {
			http.Error(w, "request failed", http.StatusUnauthorized)
			return
		}
		// Generate new token
		token, err := auth.GenerateToken(login.Email)
		if err != nil {
			http.Error(w, "request failed", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, token+"\n")
	})
}

func usersHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			http.Error(w, "authentication error", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "users\n")
	})
}

// notFoundHandler is called on other routes. It will return
// a 404 message.
func notFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404\n")
	})
}
