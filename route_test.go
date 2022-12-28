package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostalt/router"
)

func TestRouteInference(t *testing.T) {
	routes := []*router.Route{
		// TODO: Fix the stringer implementation
		// router.NewRoute([]string{http.MethodGet}, "/", testStringer),
		router.NewRoute([]string{http.MethodGet}, "/", testHandler),
		router.NewRoute([]string{http.MethodGet}, "/", testFunc),
	}

	rr := httptest.NewRecorder()

	for _, rt := range routes {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		rt.Serve(rr, req)
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
	testStringer = stringHandler{}
	testHandler  = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	})
	testFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	}
)
