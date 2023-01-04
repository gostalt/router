package router_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gostalt/router"
	"github.com/stretchr/testify/assert"
)

func TestCanRegisterCustomHandlerTransformer(t *testing.T) {
	r := router.New()

	err := r.AddHandlerTransformer(func(fn func() int) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			str := strconv.Itoa(fn())
			w.Write([]byte(str))
		})
	})
	assert.NoError(t, err)

	r.Get("/", func() int {
		return 99
	})

	server := httptest.NewServer(r)
	defer server.Close()

	assert.Equal(t, "99", get(server.URL))
}
