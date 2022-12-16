package router

import (
	"errors"
	"fmt"
	"net/http"
)

type Router struct {
	routes     []*Route
	groups     []*Group
	validators []Validator

	fallback http.Handler
}

// New creates a new router instance.
func New() *Router {
	return &Router{
		routes: make([]*Route, 0),
		groups: make([]*Group, 0),
		validators: []Validator{
			URI{}, Method{},
		},
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, err := router.findRoute(r)

	if err != nil {
		if err.Error() == "route not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 not found"))
			return
		}

		if err.Error() == "method not allowed" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 method not allowed"))
			return
		}
	}

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	match := route.Regex().FindStringSubmatch(r.RequestURI)
	fmt.Println(match)
	fmt.Println("Route params", route.params)
	for _, k := range route.params {
		i := route.Regex().SubexpIndex(k)
		r.Form.Add(k, match[i])
	}

	route.Serve(w, r)
}

func (router *Router) findRoute(r *http.Request) (*Route, error) {
	// TODO: This is very poor, but will do for now.
	for _, group := range router.groups {
		for _, route := range group.routes {
			if route.matches(router, r) {
				route.Middleware(group.middleware...)
				return route, nil
			}
		}
	}

	for _, route := range router.routes {
		if route.matches(router, r) {
			return route, nil
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

	i, ok := router.findExistingRoute(r)
	// If the route already exists, replace it.
	if ok {
		router.routes[i] = r
	} else {
		router.routes = append(router.routes, r)
	}

	return r
}

// findExistingRoute checks the router to determine if a given route already exists.
// If so, the index is returned along with `true`. Otherwise -1, false.
func (router *Router) findExistingRoute(route *Route) (int, bool) {
	for i, r := range router.routes {
		if r.path == route.path && methodsMatch(r, route) {
			return i, true
		}
	}

	return -1, false
}

func methodsMatch(routeA *Route, routeB *Route) bool {
	if len(routeA.methods) != len(routeB.methods) {
		return false
	}

	for i, v := range routeA.methods {
		if v != routeB.methods[i] {
			return false
		}
	}

	return true
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
