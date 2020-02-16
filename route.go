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

// Get defines a new `GET` route.
func Get(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodGet}, path, handler)
}

// Post defines a new `POST` route.
func Post(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodPost}, path, handler)
}

// Put defines a new `PUT` route.
func Put(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodPut}, path, handler)
}

// Patch defines a new `PATCH` route.
func Patch(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodPatch}, path, handler)
}

// Delete defines a new `DELETE` route.
func Delete(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodDelete}, path, handler)
}

// Options defines a new `OPTIONS` route.
func Options(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodOptions}, path, handler)
}

// Match defines a new route that responds to multiple http verbs.
func Match(path string, handler interface{}) *Route {
	return NewRoute([]string{http.MethodOptions}, path, handler)
}

// Any defines a new route that responds to any http verb.
func Any(path string, handler interface{}) *Route {
	verbs := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}

	return NewRoute(verbs, path, handler)
}

func Redirect(from string, to string) *Route {
	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusPermanentRedirect)
	}

	return NewRoute([]string{http.MethodGet}, from, redirect)
}
