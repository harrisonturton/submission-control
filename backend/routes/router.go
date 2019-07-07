package routes

import (
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"net/http"
)

// CreateMux creates a http.ServeMux instance, and injects
// the various dependencies into each handler.
func CreateMux(store *store.Store, logger *log.Logger) *http.ServeMux {
	var handlers = map[string]http.HandlerFunc{
		"/auth":              authHandler(store),
		"/refresh":           refreshHandler(),
		"/upload/students":   studentUploadHandler(store),
		"/upload/submission": submissionUploadHandler(store),
		"/user":              userHandler(store),
		"/assessment":        assessmentHandler(store),
		"/submissions":       submissionsHandler(store),
		"/tutorials":         tutorialHandler(store),
		"/log":               logHandler(logger),
	}
	mux := http.NewServeMux()
	for route, handler := range handlers {
		mux.Handle(route, handler)
	}
	return mux
}
