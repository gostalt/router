package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostalt/router"
	"github.com/stretchr/testify/assert"
)

func TestFindRoutesInGroups(t *testing.T) {
	r := router.New()
	r.Group(
		router.Get("test", func() string {
			return "Test"
		}),
	).Prefix("group")

	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := http.Get(server.URL + "/group/test")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGroupCanAddRoute(t *testing.T) {
	rtr := router.New()
	group := rtr.Group()

	group.Add(
		router.Get("/", func() string { return "Hello" }),
	)

	assert.Equal(t, 1, len(group.Routes()))
}
