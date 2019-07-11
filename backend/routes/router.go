package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
)

// CreateMux creates a http.ServeMux instance, and injects
// the various dependencies into each handler.
func CreateMux(store *store.Store, logger *log.Logger) *chi.Mux {
	/*var handlers = map[string]http.HandlerFunc{
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
	return mux*/
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Post("/auth", authHandler(store))
	router.Post("/refresh", refreshHandler())
	router.Post("/upload/students", studentUploadHandler(store))
	router.Post("/upload/submission", submissionUploadHandler(store))
	router.Post("/log", logHandler(logger))
	router.Get("/user", userHandler(store))
	router.Get("/assessment", assessmentHandler(store))
	router.Get("/submissions", submissionsHandler(store))
	router.Get("/tutorials", tutorialHandler(store))
	return router
}
