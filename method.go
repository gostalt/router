package router

import (
	"net/http"
)

// Method is a validator that determines whether a given Route definition matches
// the HTTP verb used in the incoming request.
type Method struct{}

func (Method) Matches(route *Route, req *http.Request) bool {
	for _, m := range route.Methods() {
		if m == req.Method {
			return true
		}
	}

	return false
}
