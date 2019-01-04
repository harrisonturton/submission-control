package routes

import (
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/auth"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
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
		// Send response
		resp := LoginResponse{token}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "request failed", http.StatusUnauthorized)
			return
		}
		w.Write(respBytes)
	})
}

// refreshHandler recieves a RefreshRequest. If the given token is valid,
// a new token is returned.
func refreshHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reject is not POST
		if r.Method != http.MethodPost {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Unmarshal POST data
		var refresh RefreshRequest
		err := json.Unmarshal(request.GetBody(r), &refresh)
		if err != nil {
			http.Error(w, "unrecognised body", http.StatusBadRequest)
			log.Println("Unrecognised body")
			return
		}
		// Verify token
		ok := auth.VerifyToken(refresh.Token)
		if !ok {
			http.Error(w, "request failed", http.StatusUnauthorized)
			log.Printf("Failed to verify login information: %v\n", err)
			return
		}
		// Get claims from the token
		claims, err := auth.ParseToken(refresh.Token)
		if err != nil {
			http.Error(w, "request failed", http.StatusUnauthorized)
			log.Println("Failed to parse the token")
			return
		}
		// Generate new token
		token, err := auth.GenerateToken(claims.Email)
		if err != nil {
			http.Error(w, "request failed", http.StatusUnauthorized)
			log.Println("Failed to generate token")
			return
		}
		w.Write([]byte(token + "\n"))
	})
}

func usersHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			http.Error(w, "authentication error", http.StatusUnauthorized)
			return
		}
		w.Write([]byte("users\n"))
	})
}

// notFoundHandler is called on other routes. It will return
// a 404 message.
func notFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404\n"))
	})
}
