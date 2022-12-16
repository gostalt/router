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
	params []param

	middleware []Middleware
}

type param struct {
	name  string
	regex *regexp.Regexp
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
	fmt.Println("route regex is", r.regex.String())
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
	if !strings.ContainsAny(path, "{}") {
		return regexp.MustCompile(path)
	}

	path = r.normalizeParamaterizedPath(path)

	rx := regexp.MustCompile("{([^}:]+):?([^}]+)?}")
	// slice of matches {}
	// first elem is full match
	// second elem is the name
	// third elem is the (optional) regex

	// dont care about first elem
	// second is appended to the route params
	res2 := rx.FindAllStringSubmatch(path, -1)
	fmt.Println("res2 is", res2)
	for _, v := range res2 {
		r.params = append(r.params, param{v[1], regexp.MustCompile(v[2])})
	}

	return regexp.MustCompile(rx.ReplaceAllString(path, "(?P<$1>$2)") + "$")
}

func (r *Route) normalizeParamaterizedPath(path string) string {
	regex := regexp.MustCompile("{([^:}]+)}")

	return regex.ReplaceAllString(path, "{$1:.+}")
}
