package router

import (
	"fmt"
	"net/http"
	"reflect"
)

// handlerTransformers is a package-level variable that contains a map between
// function signatures and the function declaration used to turn a provided
// route resolver into a valid http.Handler. This value is used at the point
// of building a handler for a newly registered route.
//
// See `handler_defaults.go` for examples of "standard" transformers.
var handlerTransformers = make(map[string]interface{})

func init() {
	for _, h := range defaultHandlers {
		if err := AddHandlerTransformer(h); err != nil {
			panic(err)
		}
	}
}

// AddHandlerTransformer adds the provided transformer function to the router,
// enabling routes to be registered that match the given signature.
//
// Note, the expected signature of `fn` is:
//
//	func(v interface{}) http.Handler
func AddHandlerTransformer(fn interface{}) error {
	t := reflect.TypeOf(fn)
	if err := validateHandlerTransformer(t); err != nil {
		return err
	}

	// Get the type of the first (and only) parameter to `fn`. This will be the
	// signature of any route resolvers that are added to the router.
	sig := t.In(0).String()
	if _, ok := handlerTransformers[sig]; ok {
		return fmt.Errorf("handler signature `%s` already exists, transformer not added", sig)
	}

	handlerTransformers[sig] = fn
	return nil
}

// buildHandler dynamically creates an http.Handler based on the function signature
// of the passed in function `fn`.
func buildHandler(v interface{}) http.Handler {
	// Retrieve the of the function from the transformer map.
	t := fmt.Sprintf("%T", v)
	f := reflect.ValueOf(handlerTransformers[t])

	// Reflect the value of `fn` and use it as the argument for the transformer
	// function. Return the value coerced to an http.Handler.
	in := []reflect.Value{reflect.ValueOf(v)}
	handler, ok := f.Call(in)[0].Interface().(http.Handler)
	if !ok {
		panic(fmt.Sprintf("expected http.Handler return type, got %T", v))
	}

	return handler
}

func validateHandlerTransformer(t reflect.Type) error {
	if t.Kind() != reflect.Func {
		return fmt.Errorf(
			"handler must be of type: func(sig interface{}) http.Handler, got: %s", t.Kind())
	}

	if t.NumIn() != 1 {
		return fmt.Errorf("handler func must have a single (1) parameter, got %d", t.NumIn())
	}

	return nil
}
