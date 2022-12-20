package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostalt/router"
)

func oneMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1"))
		next.ServeHTTP(w, r)
	})
}

func twoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("2"))
		next.ServeHTTP(w, r)
	})
}

func TestCanAddMiddlewareToRoute(t *testing.T) {
	r := router.NewRoute([]string{http.MethodGet}, "/", helloHandler).Middleware(oneMiddleware, twoMiddleware)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	r.Serve(rr, req)

	body, _ := ioutil.ReadAll(rr.Body)
	expected := "21Hello"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}

func TestCanAddMiddlewareToGroup(t *testing.T) {
	r := router.New()
	r.Group(
		router.Get("/test", helloHandler),
	).Middleware(oneMiddleware, twoMiddleware)

	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/test")

	body, _ := ioutil.ReadAll(resp.Body)
	expected := "21Hello"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}

func TestGroupMiddlewareWrapsRouteMiddleware(t *testing.T) {
	r := router.New()
	r.Group(
		router.Get("/test", helloHandler).Middleware(oneMiddleware, twoMiddleware),
	).Middleware(oneMiddleware, twoMiddleware)

	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/test")

	body, _ := ioutil.ReadAll(resp.Body)
	expected := "2121Hello"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}
