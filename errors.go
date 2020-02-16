package router

type MethodNotAllowed struct{}

func (MethodNotAllowed) Error() string {
	return "Method Not Allowed"
}

type RouteNotFound struct{}

func (RouteNotFound) Error() string {
	return "Route Not Found"
}
