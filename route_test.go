package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostalt/router"
)

func TestRouteInference(t *testing.T) {
	rtr := router.New()
	rtr.Get("/handler", testHandler)
	rtr.Get("func", testFunc)

	server := httptest.NewServer(rtr)
	defer server.Close()

	resp, err := http.Get(server.URL + "/handler")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	expected := "handler"
	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}

	resp, err = http.Get(server.URL + "/func")
	if err != nil {
		t.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	expected = "func"
	if string(body) != expected {
		t.Errorf("Got %s, wanted %s.", string(body), expected)
	}
}

type stringHandler struct{}

func (stringHandler) String() string {
	return "Test"
}

var (
	testStringer = stringHandler{}
	testHandler  = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handler"))
	})
	testFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("func"))
	}
)
