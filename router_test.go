package router_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostalt/router"
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
		// "unavailable method returns 405": {
		// 	expected: http.StatusMethodNotAllowed,
		// 	setup: func() *httptest.Server {
		// 		r := router.New()
		// 		r.Post("/", helloHandler)
		// 		return httptest.NewServer(r)
		// 	},
		// 	client: func() *http.Client {
		// 		return &http.Client{}
		// 	},
		// },
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

func TestEmptyRouteIsNotCatchall(t *testing.T) {
	router := router.New()
	router.Get("/", func() string {
		return "hello"
	})

	server := httptest.NewServer(router)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/nonexistant")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRouteDispatching(t *testing.T) {
	router := router.New()
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("test basic get", func(t *testing.T) {
		router.Get("user/profile", func() string {
			return "Hello"
		})
		assert.Equal(t, "Hello", get(server.URL+"/user/profile"))
	})

	t.Run("test basic post", func(t *testing.T) {
		router.Post("users", func() string {
			return "Hello post"
		})
		assert.Equal(t, "Hello post", post(server.URL+"/users"))
	})

	t.Run("test parameterised route", func(t *testing.T) {
		router.Get("users/{userId}/posts/{postId}", func(req *http.Request) string {
			return fmt.Sprintf("Hello %s on post %s!", req.Form.Get("userId"), req.Form.Get("postId"))
		})
		assert.Equal(t, "Hello 30 on post 28!", get(server.URL+"/users/30/posts/28"))
	})

	t.Run("test parameterised route with patterns", func(t *testing.T) {
		router.Get("posts/{postId:[0-9]+}", func(req *http.Request) string {
			return fmt.Sprintf("Hello post %s!", req.Form.Get("postId"))
		})
		assert.Equal(t, "Hello post 28!", get(server.URL+"/posts/28"))
	})

	t.Run("duplicate records uses last registered", func(t *testing.T) {
		router.Get("duplicate", func() string {
			return "first"
		})

		router.Get("duplicate", func() string {
			return "second"
		})

		assert.Equal(t, "second", get(server.URL+"/duplicate"))
	})
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
