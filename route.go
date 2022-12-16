package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Route is a single entrypoint into the router.
type Route struct {
	methods []string
	path    string
	handler http.Handler
	regex   *regexp.Regexp

	// The {} bits of a route
	params []string

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
	return newHandlerRoute(methods, path, buildHandler(handler))
}

func newHandlerRoute(methods []string, path string, handler http.Handler) *Route {
	if path[0] != '/' {
		path = "/" + path
	}

	r := &Route{
		methods: methods,
		path:    path,
		handler: handler,
	}

	r.regex = r.calculateRouteRegex(path)
	fmt.Println("Route regex is", r.regex.String())
	return r
}

func (route *Route) Serve(w http.ResponseWriter, r *http.Request) {
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
func Match(verbs []string, path string, handler interface{}) *Route {
	return NewRoute(verbs, path, handler)
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

// matches determines if the route matches the incoming request.
func (r *Route) matches(router *Router, req *http.Request) bool {
	for _, v := range router.validators {
		if !v.Matches(r, req) {
			return false
		}
	}

	return true
}

func (r *Route) Methods() []string {
	return r.methods
}

func (r *Route) Regex() *regexp.Regexp {
	return r.regex
}

func (r *Route) calculateRouteRegex(path string) *regexp.Regexp {
	rx := regexp.MustCompile("{([^}]+)}")
	r.params = r.getParamsFromURI(path)
	fmt.Println("The params are", r.params)
	return regexp.MustCompile(rx.ReplaceAllString(path, "(?P<$1>.+)") + "$")
}

func (r *Route) getParamsFromURI(uri string) []string {
	var params []string
	rx := regexp.MustCompile("{([^}]+)}")
	res := rx.FindAllString(uri, -1)
	fmt.Println("The res is", res)
	for _, v := range res {
		fmt.Println("param", v)
		v = strings.TrimPrefix(v, "{")
		v = strings.TrimSuffix(v, "}")
		params = append(params, v)
	}

	return params
}
