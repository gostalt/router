package router

import (
	"net/http"
)

type Method struct{}

func (Method) Matches(route *Route, req *http.Request) bool {
	for _, m := range route.Methods() {
		if m == req.Method {
			return true
		}
	}

	return false
}
