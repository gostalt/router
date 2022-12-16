package router

import (
	"net/http"
)

type Validator interface {
	Matches(*Route, *http.Request) bool
}
