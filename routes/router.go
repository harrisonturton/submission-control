package routes

import (
	"github.com/harrisonturton/submission-control/store"
	"net/http"
)

// CreateMux creates a http.ServeMux instance, and injects
// the various dependencies into each handler.
func CreateMux(store *store.Store) *http.ServeMux {
	var handlers = map[string]http.HandlerFunc{
		"/auth":  authHandler(store),
		"/users": usersHandler(store),
		"/":      notFoundHandler(),
	}
	mux := http.NewServeMux()
	for route, handler := range handlers {
		mux.HandleFunc(route, handler)
	}
	return mux
}
