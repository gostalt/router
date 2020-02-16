package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestCanAddMiddleware(t *testing.T) {
	r := NewRoute(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})).Middleware(oneMiddleware, twoMiddleware)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	r.serve(rr, req)

	body, _ := ioutil.ReadAll(rr.Body)
	expected := "21Hello"

	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}
