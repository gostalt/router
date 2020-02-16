package router

import (
	"fmt"
	"net/http"
)

// Route is a single entrypoint into the router.
type Route struct {
	method  string
	path    string
	handler http.Handler

	middleware []Middleware
}

// Middleware defines additional logic on a single route definition by wrapping the
// route's handler in extra layers of logic.
func (route *Route) Middleware(middleware ...Middleware) *Route {
	route.middleware = append(route.middleware, middleware...)
	return route
}

// NewRoute creates a new route definition for a given method, path and handler.
func NewRoute(method string, path string, handler interface{}) *Route {
	switch handler.(type) {
	case string:
		return newStringRoute(method, path, handler.(string))
	case fmt.Stringer:
		return newStringerRoute(method, path, handler.(fmt.Stringer))
	case func(http.ResponseWriter, *http.Request):
		return newHandlerRoute(method, path, http.HandlerFunc(handler.(func(http.ResponseWriter, *http.Request))))
	default:
		return newHandlerRoute(method, path, handler.(http.Handler))
	}
}

func newHandlerRoute(method string, path string, handler http.Handler) *Route {
	return &Route{
		method:  method,
		path:    path,
		handler: handler,
	}
}

func newStringerRoute(method string, path string, handler fmt.Stringer) *Route {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(handler.String()))
	})
	return newHandlerRoute(method, path, f)
}

func newStringRoute(method string, path string, response string) *Route {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	})

	return newHandlerRoute(method, path, h)
}

func (route *Route) serve(w http.ResponseWriter, r *http.Request) {
	handler := route.handler
	for _, m := range route.middleware {
		handler = m(handler)
	}

	handler.ServeHTTP(w, r)
}
