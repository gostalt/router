package router_test

import (
	"net/http"
	"net/http/httptest"
	"router"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRoutesInGroups(t *testing.T) {
	r := router.New()
	r.Group(
		router.Get("/test", "Test"),
	).Prefix("/group")

	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/group/test")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGroupCanAddRoute(t *testing.T) {
	group := router.NewGroup()

	group.Add(
		router.NewRoute([]string{http.MethodGet}, "/", func() string { return "Hello" }),
	)

	assert.Equal(t, 1, len(group.Routes()))
}
