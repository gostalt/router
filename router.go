package router

import (
	"errors"
	"fmt"
	"net/http"
)

type Router struct {
	groups       []*Group
	defaultGroup *Group

	validators []Validator

	fallback http.Handler
	// middleware are handlers that wrap all route definitions on this router instance.
	middleware []Middleware

	transformers map[string]interface{}
}

// New creates a new Router instance.
func New() *Router {
	rtr := &Router{
		validators: []Validator{
			URI{}, Method{},
		},
		transformers: map[string]interface{}{},
	}

	def := newGroup(rtr)
	rtr.defaultGroup = def
	rtr.groups = []*Group{def}

	for _, h := range defaultHandlers {
		if err := rtr.AddHandlerTransformer(h); err != nil {
			panic(err)
		}
	}

	fmt.Println("There are x transformers", len(rtr.transformers))

	return rtr
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
	for _, k := range route.params {
		i := route.Regex().SubexpIndex(k)
		r.Form.Add(k, match[i])
	}

	route.Serve(w, r)
}

func (router *Router) findRoute(r *http.Request) (*Route, error) {
	for _, group := range router.groups {
		for _, route := range group.routes {
			if route.matches(router, r) {
				return route, nil
			}
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

// Any creates a new route definition that responds to any HTTP verb.
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

// Redirect creates a new route definition that redirects from the `from` URI to
// the `to` URI. The redirect uses the Permanent Redirect status code 308.
func (router *Router) Redirect(from string, to string) *Route {
	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusPermanentRedirect)
	}

	return router.addRoute([]string{http.MethodGet}, from, redirect)
}

func (router *Router) addRoute(methods []string, path string, handler interface{}) *Route {
	r := NewRoute(methods, path, handler)

	router.defaultGroup.Add(r)

	r.router = router

	return r
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

// Group creates a new route Group for the Router instance.
func (router *Router) Group(routes ...*Route) *Group {
	g := newGroup(router, routes...)

	for _, route := range routes {
		route.router = router
	}

	router.groups = append(router.groups, g)

	return g
}

// Fallback defines a "default" route for the Router instance. If a visited URI
// does not have a corresponding route definition, the Fallback handler is
// called for the request.
func (router *Router) Fallback(handler http.Handler) *Router {
	router.fallback = handler
	return router
}

// Middleware appends the given middleware `fns` to the Router instance.
func (router *Router) Middleware(fns ...Middleware) *Router {
	router.middleware = append(router.middleware, fns...)
	return router
}
