package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/harrisonturton/submission-control/backend/ci"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
)

// CreateMux creates a http.ServeMux instance, and injects
// the various dependencies into each handler.
func CreateMux(store *store.Store, logger *log.Logger, ci *ci.Ci) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Post("/auth", authHandler(store))
	router.Post("/refresh", refreshHandler())
	router.Post("/upload/students", studentUploadHandler(store))
	router.Post("/upload/submission/feedback/{submissionID}", submissionFeedbackHandler(store))
	router.Post("/upload/submission", submissionUploadHandler(store, ci))
	router.Post("/log", logHandler(logger))
	router.Get("/user", userHandler(store))
	router.Get("/assessment", assessmentHandler(store))
	router.Get("/submissions", submissionsHandler(store))
	router.Get("/submission/file/{submissionID}", submissionFileHandler(store))
	router.Get("/tutorials", tutorialHandler(store))
	return router
}
