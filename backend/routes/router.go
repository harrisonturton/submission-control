package routes

import (
	"github.com/harrisonturton/submission-control/backend/store"
	"net/http"
)

// CreateMux creates a http.ServeMux instance, and injects
// the various dependencies into each handler.
func CreateMux(store *store.Store) *http.ServeMux {
	var handlers = map[string]http.HandlerFunc{
		"/auth":            authHandler(store),
		"/refresh":         refreshHandler(),
		"/upload/students": studentUploadHandler(store),
		"/user":            userHandler(store),
		"/assessment":      assessmentHandler(store),
		"/submissions":     submissionsHandler(store),
		"/tutorials":       tutorialHandler(store),
	}
	mux := http.NewServeMux()
	for route, handler := range handlers {
		mux.Handle(route, handler)
	}
	return mux
}
