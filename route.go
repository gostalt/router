package router

import (
	"net/http"
)

type Route struct {
	path    string
	handler http.Handler

	middleware []Middleware
}

func (r *Route) Middleware(middleware ...Middleware) *Route {
	r.middleware = append(r.middleware, middleware...)
	return r
}

func NewRoute(path string, handler http.Handler) *Route {
	return &Route{
		path:    path,
		handler: handler,
	}
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := route.handler
	for _, m := range route.middleware {
		handler = m(handler)
	}

	handler.ServeHTTP(w, r)
}
