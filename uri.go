package router

import (
	"fmt"
	"net/http"
)

type URI struct{}

func (URI) Matches(route *Route, req *http.Request) bool {
	fmt.Println(route.Regex().String())
	fmt.Println(req.RequestURI)
	fmt.Printf("%s matches %s? %v\n", route.Regex().String(), req.RequestURI, route.Regex().MatchString(req.RequestURI))
	return route.Regex().MatchString(req.RequestURI)
}
