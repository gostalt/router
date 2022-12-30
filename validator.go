package router

import (
	"net/http"
)

// Validator determines whether an incoming request matches a route definition.
type Validator interface {
	Matches(*Route, *http.Request) bool
}
