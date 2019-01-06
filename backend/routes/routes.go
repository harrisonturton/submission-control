package routes

import (
	"encoding/json"
	"fmt"
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

const (
	// These are the messages given for various error routes
	errorNotFoundMessage     = "not found"
	errorUnauthorizedMessage = "unauthorized"
)

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
			badRequestHandler("unrecognised body").ServeHTTP(w, r)
			return
		}
		// Verify login information
		ok, err := auth.Authenticate(store, login.Email, login.Password)
		if !ok || err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Generate new token
		token, err := auth.GenerateToken(login.Email)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := TokenResponse{
			StatusCode: http.StatusOK,
			Token:      token,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		log.Println(string(respBytes))
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
			badRequestHandler("unrecognised body").ServeHTTP(w, r)
			return
		}
		// Verify token
		ok := auth.VerifyToken(refresh.Token)
		if !ok {
			log.Printf("Failed to verify login information: %v\n", err)
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Get claims from the token
		claims, err := auth.ParseToken(refresh.Token)
		if err != nil {
			log.Println("Failed to parse the token")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Generate new token
		token, err := auth.GenerateToken(claims.Email)
		if err != nil {
			log.Println("Failed to generate token")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := TokenResponse{
			StatusCode: http.StatusOK,
			Token:      token,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		w.Write(respBytes)
	})
}

func usersHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		w.Write([]byte("users\n"))
	})
}

func addPreflightHeaders(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			log.Println("Handling OPTIONS request")
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// notFoundHandler is called on other routes. It will return
// a 404 message.
func notFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(
			`{"status":%d,"message":"%s"}`, http.StatusNotFound, errorNotFoundMessage)
		http.Error(w, resp, http.StatusNotFound)
	})
}

// unauthorizedHandler responds with a status unauthorized message
func unauthorizedHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(
			`{"status":%d,"message":"%s"}`, http.StatusUnauthorized, errorUnauthorizedMessage)
		http.Error(w, resp, http.StatusUnauthorized)
	})
}

// badRequestHandler responds with a bad request code and a custom message
func badRequestHandler(message string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(
			`{"status":%d,"message":"%s"}`, http.StatusBadRequest, message)
		http.Error(w, resp, http.StatusBadRequest)
	})
}
