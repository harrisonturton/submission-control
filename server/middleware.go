package server

import (
	"fmt"
	"github.com/harrisonturton/submission-control/request"
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
// in the request.Context object. See server/context.go for more details.
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

// logAll will log information about all incoming requests. It expects the
// request body to be already read & stored in the request context.
func logAll(log *log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request := displayRequest(r)
			log.Println(request)
			h.ServeHTTP(w, r)
		})
	}
}

// displayRequest returns a pretty-printed summary of a http Request.
// It expects the request body to be already read and put into the request
// context.
func displayRequest(r *http.Request) (result string) {
	if r.URL != nil {
		result += fmt.Sprintf("%s %s\n", r.Method, r.URL.Path)
	} else {
		result += fmt.Sprintf("%s [no url]\n", r.Method)
	}
	result += fmt.Sprintf("Content Length: %d\n", r.ContentLength)
	body := request.GetBody(r)
	result += string(body) + "\n"
	return
}