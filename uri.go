package router

import (
	"net/http"
)

// URI is a Validator that determines whether a given Route definition matches
// the incoming request URI.
type URI struct{}

func (URI) Matches(route *Route, req *http.Request) bool {
	return route.Regex().MatchString(req.RequestURI)
}
