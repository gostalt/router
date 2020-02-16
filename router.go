package router

import (
	"errors"
	"net/http"
)

// Router is an http.Handler that you can register routes against. It can be passed
// to a http.Server.
type Router struct {
	routes []*Route
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// determine the route and execute it here.
	route, err := router.findRoute(r.RequestURI, r.Method)

	if errors.Is(err, MethodNotAllowed{}) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if errors.Is(err, RouteNotFound{}) {
		http.Error(w, "404 — Route not Found", 404)
		return
	}

	route.serve(w, r)
}

func (router *Router) findRoute(path string, method string) (*Route, error) {
	for _, route := range router.routes {
		if route.path == path {
			if route.method != method {
				return &Route{}, MethodNotAllowed{}
			}

			return route, nil
		}
	}

	return &Route{}, RouteNotFound{}
}

// Get defines a new `GET` route on the router, at the given path.
func (router *Router) Get(path string, handler interface{}) *Route {
	r := NewRoute(http.MethodGet, path, handler)
	router.routes = append(router.routes, r)

	return r
}

// New creates a new router instance.
func New() *Router {
	return &Router{
		routes: []*Route{},
	}
}
