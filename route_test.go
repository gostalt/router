package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteInference(t *testing.T) {
	routes := []*Route{
		NewRoute(http.MethodGet, "/string", "Test"),
		NewRoute(http.MethodGet, "/stringer", stringHandler{}),
		NewRoute(http.MethodGet, "/handler", testHandler),
		NewRoute(http.MethodGet, "/handler2", testFunc),
	}

	rr := httptest.NewRecorder()

	for _, rt := range routes {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		rt.serve(rr, req)
		body, _ := ioutil.ReadAll(rr.Body)
		expected := "Test"

		if string(body) != expected {
			t.Errorf("Got %s, wanted %s.", string(body), expected)
		}
	}
}

type stringHandler struct{}

func (stringHandler) String() string {
	return "Test"
}

var (
	testString   = "Test"
	testStringer = stringHandler{}
	testHandler  = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	})
	testFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	}
)
