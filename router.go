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
			for _, m := range route.methods {
				if m == method {
					return route, nil
				}
			}

			return &Route{}, MethodNotAllowed{}
		}
	}

	return &Route{}, RouteNotFound{}
}

// Get defines a new `GET` route on the router, at the given path.
func (router *Router) Get(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodGet}, path, handler)
}

// Post defines a new `POST` route on the router, at the given path.
func (router *Router) Post(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodPost}, path, handler)
}

// Put defines a new `PUT` route on the router, at the given path.
func (router *Router) Put(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodPut}, path, handler)
}

// Patch defines a new `PATCH` route on the router, at the given path.
func (router *Router) Patch(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodPatch}, path, handler)
}

// Delete defines a new `DELETE` route on the router, at the given path.
func (router *Router) Delete(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodDelete}, path, handler)
}

// Options defines a new `OPTIONS` route on the router, at the given path.
func (router *Router) Options(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodOptions}, path, handler)
}

// Match defines a new route that responds to multiple http verbs.
func (router *Router) Match(path string, handler interface{}) *Route {
	return router.addRoute([]string{http.MethodOptions}, path, handler)
}

// Any defines a new route that responds to any http verb.
func (router *Router) Any(path string, handler interface{}) *Route {
	verbs := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}

	return router.addRoute(verbs, path, handler)
}

func (router *Router) Redirect(from string, to string) *Route {
	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusPermanentRedirect)
	}

	return router.addRoute([]string{http.MethodGet}, from, redirect)
}

func (router *Router) addRoute(methods []string, path string, handler interface{}) *Route {
	r := NewRoute(methods, path, handler)
	router.routes = append(router.routes, r)

	return r
}

// New creates a new router instance.
func New() *Router {
	return &Router{
		routes: []*Route{},
	}
}
