package router

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	mux.HandleFunc("/users", router.authMiddleware(router.usersHandler))
	mux.HandleFunc("/", router.authMiddleware(router.notFoundHandler))
	router.mux = mux
	router.logger = logger
	router.store = store
	return router
}

func (router *Router) authMiddleware(nextHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, ok := r.Header["Authorization"]
		if !ok || len(tokenStr) == 0 {
			http.Error(w, "authentication failure", http.StatusUnauthorized)
			return
		}
		var claims Claims
		token, err := jwt.ParseWithClaims(tokenStr[0], &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "authentication failure", http.StatusUnauthorized)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return SigningKey, nil
		})
		if err != nil {
			router.logger.Println("claim: " + tokenStr[0])
			router.logger.Printf("Failed to parse claims: %v\n", err)
			http.Error(w, "authentication failure", http.StatusUnauthorized)
			return
		}
		_, ok = token.Claims.(*Claims)
		if !ok || !token.Valid {
			http.Error(w, "authentication failure", http.StatusUnauthorized)
			return
		}
		nextHandler(w, r)
	}
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
