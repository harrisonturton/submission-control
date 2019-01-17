package routes

import (
	"github.com/harrisonturton/submission-control/backend/store"
	"net/http"
)

// CreateMux creates a http.ServeMux instance, and injects
// the various dependencies into each handler.
func CreateMux(store *store.Store) *http.ServeMux {
	var handlers = map[string]http.HandlerFunc{
		"/auth":        authHandler(store),
		"/refresh":     refreshHandler(store),
		"/user":        userHandler(store),
		"/enrol":       enrolHandler(store),
		"/assessment":  assessmentHandler(store),
		"/submissions": submissionHandler(store),
		"/":            notFoundHandler(),
	}
	mux := http.NewServeMux()
	for route, handler := range handlers {
		mux.Handle(route, handler)
	}
	return mux
}
