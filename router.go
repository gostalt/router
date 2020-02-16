package router

import (
	"errors"
	"net/http"
)

type Router struct {
	routes []*Route
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// determine the route and execute it here.
	route, err := router.findRouteForPath(r.RequestURI)
	if err != nil {
		http.Error(w, "404 — Route not Found", 404)
		return
	}
	route.ServeHTTP(w, r)
}

func (router *Router) findRouteForPath(path string) (*Route, error) {
	for _, route := range router.routes {
		if route.path == path {
			return route, nil
		}
	}

	return &Route{}, errors.New("Route not found")
}

func (router *Router) Get(path string, handler http.Handler) *Route {
	route := &Route{
		path:    path,
		handler: handler,
	}
	router.routes = append(router.routes, route)

	return route
}

func New() *Router {
	return &Router{
		routes: []*Route{},
	}
}
