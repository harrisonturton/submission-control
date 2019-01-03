package routes

import (
	"encoding/json"
	"fmt"
	"github.com/harrisonturton/submission-control/auth"
	"github.com/harrisonturton/submission-control/request"
	"github.com/harrisonturton/submission-control/store"
	"net/http"
)

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

func notFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404\n")
	})
}
