package router

import (
	"fmt"
	"net/http"
	"reflect"
)

var handlerTransformers map[string]interface{} = map[string]interface{}{
	"func() string":                            new__func_ret_string__Handler,
	"func(*http.Request) string":               new__func_http_Request_ret_string__Handler,
	"fmt.Stringer":                             new__Stringer__Handler,
	"http.HandlerFunc":                         new__http_HandlerFunc__Handler,
	"func(http.ResponseWriter, *http.Request)": new__func_http_ResponseWriter_http_Request__Handler,
	"http.Handler":                             new_http_Handler__Handler,
}

// TODO: Would really rather get the sig dynamically
func AddHandlerTransformer(sig string, fn interface{}) error {
	if _, ok := handlerTransformers[sig]; ok {
		// Handler signature already exists, return an error.
		return fmt.Errorf("handler signature `%s` already exists, transformer not added", sig)
	}

	handlerTransformers[sig] = fn
	return nil
}

func new_http_Handler__Handler(fn http.Handler) http.Handler {
	return fn
}

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

func new__func_http_ResponseWriter_http_Request__Handler(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(fn)
}

// buildHandler dynamically creates an http.Handler based on the function signature
// of the passed in function `fn`.
func buildHandler(fn interface{}) http.Handler {
	// Retrieve the signature of the function. This is the key in the transformer map.
	a := fmt.Sprintf("%T", fn)
	handler := reflect.ValueOf(handlerTransformers[a])

	// Reflect the value of `fn` and use it as the argument for the transformer
	// function. Return the value coerced to an http.Handler.
	// TODO: Ensure this is safe to do, probably at the point of registration.
	// Could also `val, ok :=` it here and panic?
	in := []reflect.Value{reflect.ValueOf(fn)}
	return handler.Call(in)[0].Interface().(http.Handler)
}

func makeFailedHandler(handler interface{}) http.HandlerFunc {
	msg := fmt.Sprintf("Unable to create handler for type %T", handler)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	}
}
