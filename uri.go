package router

import (
	"net/http"
)

type URI struct{}

func (URI) Matches(route *Route, req *http.Request) bool {
	return route.Regex().MatchString(req.RequestURI)
}
