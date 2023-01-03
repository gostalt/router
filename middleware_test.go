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

func threeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("3"))
		next.ServeHTTP(w, r)
	})
}

func TestCanAddMiddlewareToRoute(t *testing.T) {
	rtr := router.New()
	rtr.Get("/", helloHandler).Middleware(oneMiddleware, twoMiddleware)

	server := httptest.NewServer(rtr)
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
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

func TestMiddlewareExecutionOrder(t *testing.T) {
	r := router.New()
	r.Middleware(oneMiddleware)

	r.Group(
		router.Get("group-middleware", func() string {
			return "middleware"
		}).Middleware(threeMiddleware),
	).Middleware(twoMiddleware)

	r.Get("route-middleware", func() string {
		return "middleware"
	}).Middleware(twoMiddleware)

	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/group-middleware")

	body, _ := ioutil.ReadAll(resp.Body)
	expected := "123middleware"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}

	resp, _ = http.Get(server.URL + "/route-middleware")

	body, _ = ioutil.ReadAll(resp.Body)
	expected = "12middleware"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}

func TestMiddlewareDoesntDuplicate(t *testing.T) {
	r := router.New()
	r.Middleware(oneMiddleware)
	r.Get("middleware", func() string {
		return "middleware"
	}).Middleware(twoMiddleware)

	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/middleware")
	resp, _ = http.Get(server.URL + "/middleware")
	resp, _ = http.Get(server.URL + "/middleware")

	body, _ := ioutil.ReadAll(resp.Body)
	expected := "12middleware"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}
