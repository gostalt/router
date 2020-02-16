package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Routers must respond with a 404 status code if the path
// given cannot be found.
func TestUnregisteredRoutesReturn404(t *testing.T) {
	server := httptest.NewServer(New())
	defer server.Close()

	resp, _ := http.Get(server.URL + "/404")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestRegisteredRouteCanBeAccessed(t *testing.T) {
	r := New()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))
	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL)

	body, _ := ioutil.ReadAll(resp.Body)
	expected := "Hello"
	if string(body) != expected {
		t.Errorf("Expected `%s`, got `%s`", expected, string(body))
	}
}
