package router

import (
	"fmt"
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

func NewRoute(path string, handler interface{}) *Route {
	switch handler.(type) {
	case string:
		return newStringRoute(path, handler.(string))
	case fmt.Stringer:
		return newStringerRoute(path, handler.(fmt.Stringer))
	case func(http.ResponseWriter, *http.Request):
		return newHandlerRoute(path, http.HandlerFunc(handler.(func(http.ResponseWriter, *http.Request))))
	default:
		return newHandlerRoute(path, handler.(http.Handler))
	}
}

func newHandlerRoute(path string, handler http.Handler) *Route {
	return &Route{
		path:    path,
		handler: handler,
	}
}

func newStringerRoute(path string, handler fmt.Stringer) *Route {
	return newStringRoute(path, handler.String())
}

func newStringRoute(path string, response string) *Route {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	})

	return newHandlerRoute(path, h)
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := route.handler
	for _, m := range route.middleware {
		handler = m(handler)
	}

	handler.ServeHTTP(w, r)
}
