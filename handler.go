package router

import (
	"fmt"
	"net/http"
)

// new__func_ret_string__Handler returns a valid http.Handler from a function that
// matches the following signature:
//
//	func() string
func new__func_ret_string__Handler(fn func() string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fn()))
	})
}

// new__func_http_Request_ret_string__Handler returns a valid http.Handler from a
// function that matches the following signature:
//
//	func(*http.Request) string
func new__func_http_Request_ret_string__Handler(fn func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fn(r)))
	})
}

// new__Stringer__Handler returns a valid http.Handler from a value that satisfies
// the fmt.Stringer interface, that is:
//
//	type Stringer interface {
//		String() string
//	}
func new__Stringer__Handler(fn fmt.Stringer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fn.String()))
	})
}

func new__http_HandlerFunc__Handler(fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(fn)
}

func buildHandler(fn interface{}) http.Handler {
	switch v := fn.(type) {
	case http.HandlerFunc:
		return new__http_HandlerFunc__Handler(v)
	case http.Handler:
		return v
	case func() string:
		return new__func_ret_string__Handler(v)
	case func(*http.Request) string:
		return new__func_http_Request_ret_string__Handler(v)
	case fmt.Stringer:
		return new__Stringer__Handler(v)
	default:
		return makeFailedHandler(v)
	}
}

func makeFailedHandler(handler interface{}) http.HandlerFunc {
	msg := fmt.Sprintf("Unable to create handler for type %T", handler)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	}
}
