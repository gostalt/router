package router

import (
	"net/http"
)

// defaultHandlers define the preconfigured transformers for the Router.
var defaultHandlers = []interface{}{
	func(fn func() string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fn()))
		})
	},
	func(fn func(*http.Request) string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fn(r)))
		})
	},
	func(fn http.HandlerFunc) http.Handler {
		return http.HandlerFunc(fn)
	},
	func(fn func(http.ResponseWriter, *http.Request)) http.Handler {
		return http.HandlerFunc(fn)
	},
	func(fn http.Handler) http.Handler {
		return fn
	},
}
