package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindRoutesInGroups(t *testing.T) {
	r := New()
	rt := NewRoute([]string{http.MethodGet}, "/test", "Test")
	r.Group(rt).Prefix("/group")

	server := httptest.NewServer(r)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/group/test")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
