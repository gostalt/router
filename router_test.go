package router_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"router"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helloHandler is a simple handler used in tests.
func helloHandler() string { return "Hello" }

func TestResponseCodes(t *testing.T) {
	type testcase struct {
		expected int
		setup    func() *httptest.Server
		client   func() *http.Client
	}

	cases := map[string]testcase{
		"unregistered routes return 404": {
			expected: http.StatusNotFound,
			setup: func() *httptest.Server {
				return httptest.NewServer(router.New())
			},
			client: func() *http.Client {
				return &http.Client{}
			},
		},
		"unavailable method returns 405": {
			expected: http.StatusMethodNotAllowed,
			setup: func() *httptest.Server {
				r := router.New()
				r.Post("/", helloHandler)
				return httptest.NewServer(r)
			},
			client: func() *http.Client {
				return &http.Client{}
			},
		},
		"redirect routes return 308": {
			expected: http.StatusPermanentRedirect,
			setup: func() *httptest.Server {
				r := router.New()
				r.Redirect("/", "/new")
				return httptest.NewServer(r)
			},
			client: func() *http.Client {
				return &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := tc.setup()
			defer server.Close()

			resp, _ := tc.client().Get(server.URL)

			assert.Equal(t, tc.expected, resp.StatusCode)
		})
	}
}

func TestRouteDispatching(t *testing.T) {
	router := router.New()
	server := httptest.NewServer(router)
	defer server.Close()

	router.Get("user/profile", func() string {
		return "Hello"
	})
	assert.Equal(t, "Hello", get(server.URL+"/user/profile"))

	router.Post("users", func() string {
		return "Hello post"
	})
	assert.Equal(t, "Hello post", post(server.URL+"/users"))

	router.Get("users/.+", func(req *http.Request) string {
		return fmt.Sprintf("Hello %s!", req.Form.Get("id"))
	})
	assert.Equal(t, "Hello 30!", get(server.URL+"/users/30"))
}

// get is a convenience method that fires off a GET request and assumes a positive
// response with no errors. If errors occur, a panic is thrown.
func get(uri string) string {
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

// post is a convenience method that fires off a POST request and assumes a positive
// response with no errors. If errors occur, a panic is thrown.
func post(uri string) string {
	resp, err := http.Post(uri, "text/plain", nil)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}
