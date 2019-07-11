package server

import (
	"fmt"
	"github.com/harrisonturton/submission-control/backend/request"
	"log"
	"net/http"
)

// Middleware functions are ones that wrap an existing http.Handler,
// adding extra functionality. They "intercept" requests before they're
// routed to their handler.
// Multiple middleware functions can be wrapped around a handler using
// addMiddleware().
type Middleware func(http.Handler) http.Handler

// addMiddleware wraps multiple Middleware functions around a single
// handler.
func addMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// Middleware functions

// attachContext will verify the request, read the body, and store these results
// in the request.Context object. See request/request.go for more details.
func attachContext() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = request.NewContextWithAuth(ctx, r)
			ctx = request.NewContextWithBody(ctx, r)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// addPreflightHeaders will add the required headers to allow CORS requests.
// It adds the same headers to all responses.
// If a HTTP OPTIONS request comes through, it will finish here.
func addPreflightHeaders() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "content-type, Content-Type, Origin, Authorization")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
			} else {
				h.ServeHTTP(w, r)
			}
		})
	}
}

// logAll will log information about all incoming requests. It expects the
// request body to be already read & stored in the request context.
func logAll(log *log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request := displayRequest(r)
			log.Printf(request)
			h.ServeHTTP(w, r)
		})
	}
}

// displayRequest returns a pretty-printed summary of a http Request.
// It expects the request body to be already read and put into the request
// context.
func displayRequest(r *http.Request) (result string) {
	result += fmt.Sprintf("Recieved %s request\n", r.Method)
	result += fmt.Sprintf("%s %s\n", r.Method, r.RequestURI)
	result += fmt.Sprintf("From: %s\n", r.RemoteAddr)
	if contentType, ok := r.Header["Content-Type"]; ok {
		result += fmt.Sprintf("Content-Type: %s\n", contentType)
	}
	if r.ContentLength > 0 {
		result += fmt.Sprintf("Content Length: %d\n", r.ContentLength)
		body := request.GetBody(r)
		result += string(body)
	}
	return
}
