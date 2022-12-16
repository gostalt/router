package router

import (
	"errors"
	"net/http"
)

type Router struct {
	routes []*Route
	groups []*Group

	fallback http.Handler
}

// New creates a new router instance.
func New() *Router {
	return &Router{
		routes: make([]*Route, 0),
		groups: make([]*Group, 0),
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, _ := router.findRoute(r)

	route.serve(w, r)
}

func (router *Router) findRoute(r *http.Request) (*Route, error) {
	path := r.RequestURI
	method := r.Method
	// TODO: This is very poor, but will do for now.
	for _, group := range router.groups {
		for _, route := range group.routes {
			if group.prefix+route.path == path {
				for _, m := range route.methods {
					if m == method {
						route.Middleware(group.middleware...)
						return route, nil
					}
				}

				return &Route{}, errors.New("method not allowed")
			}
		}
	}

	for _, route := range router.routes {
		if route.path == path {
			for _, m := range route.methods {
				if m == method {
					return route, nil
				}
			}

			return &Route{}, errors.New("method not allowed")
		}
	}

	return &Route{}, errors.New("route not found")
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
func (router *Router) Match(verbs []string, path string, handler interface{}) *Route {
	return router.addRoute(verbs, path, handler)
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

func (router *Router) Group(routes ...*Route) *Group {
	g := NewGroup(routes...)

	router.groups = append(router.groups, g)

	return g
}

func (router *Router) Fallback(handler interface{}) {
	response := handler.(string)

	router.fallback = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	})
}

/* Route definition maybe should have a slice of "matchers" - this is where the actual route URI
would be, as well as methods, etc, etc. Eg gorilla has ability to match on header, scheme, etc.
*/
