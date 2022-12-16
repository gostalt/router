package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"router"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helloHandler is a simple handler used in tests.
func helloHandler() string { return "Hello" }

func TestUnregisteredRoutesReturn404(t *testing.T) {
	server := httptest.NewServer(router.New())
	defer server.Close()

	resp, _ := http.Get(server.URL + "/404")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRegisteredRouteCanBeAccessed(t *testing.T) {
	r := router.New()
	r.Get("/", helloHandler)
	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("Hello"), body)
}

func TestRoutesCanOnlyBeAccessedByRegisteredMethods(t *testing.T) {
	r := router.New()
	r.Get("/", helloHandler)
	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Post(server.URL, "application/json", nil)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestRedirectRoute(t *testing.T) {
	r := router.New()
	r.Redirect("/", "/new")
	server := httptest.NewServer(r)
	defer server.Close()

	// Prevent redirects from occuring, so we can check the status code of the request.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, _ := client.Get(server.URL)
	assert.Equal(t, http.StatusPermanentRedirect, resp.StatusCode)
}
