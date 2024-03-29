package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestRoutesCanOnlyBeAccessedByRegisteredMethods(t *testing.T) {
	r := New()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))
	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Post(server.URL, "application/json", nil)

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestRedirectRoute(t *testing.T) {
	r := New()
	r.Redirect("/", "/new")
	server := httptest.NewServer(r)
	defer server.Close()

	// Prevent redirects from occuring, so we can check the status code of the request.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, _ := client.Get(server.URL)

	if resp.StatusCode != http.StatusPermanentRedirect {
		t.Errorf("Expected %d, got %d", http.StatusPermanentRedirect, resp.StatusCode)
	}
}
