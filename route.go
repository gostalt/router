package router

import (
	"fmt"
	"net/http"
)

// Route is a single entrypoint into the router.
type Route struct {
	methods []string
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
func NewRoute(methods []string, path string, handler interface{}) *Route {
	switch handler.(type) {
	case string:
		return newStringRoute(methods, path, handler.(string))
	case fmt.Stringer:
		return newStringerRoute(methods, path, handler.(fmt.Stringer))
	case func(http.ResponseWriter, *http.Request):
		return newHandlerRoute(methods, path, http.HandlerFunc(handler.(func(http.ResponseWriter, *http.Request))))
	default:
		return newHandlerRoute(methods, path, handler.(http.Handler))
	}
}

func newHandlerRoute(methods []string, path string, handler http.Handler) *Route {
	return &Route{
		methods: methods,
		path:    path,
		handler: handler,
	}
}

func newStringerRoute(methods []string, path string, handler fmt.Stringer) *Route {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(handler.String()))
	})
	return newHandlerRoute(methods, path, f)
}

func newStringRoute(methods []string, path string, response string) *Route {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	})

	return newHandlerRoute(methods, path, h)
}

func (route *Route) serve(w http.ResponseWriter, r *http.Request) {
	handler := route.handler
	for _, m := range route.middleware {
		handler = m(handler)
	}

	handler.ServeHTTP(w, r)
}
