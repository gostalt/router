package router

// MethodNotAllowed denotes that finding a route definition was successful, but
// the method used to fire it is not listed on the route definition.
type MethodNotAllowed struct{}

func (MethodNotAllowed) Error() string {
	return "Method Not Allowed"
}

// RouteNotFound denotes that a route at the given path was not found on the router.
type RouteNotFound struct{}

func (RouteNotFound) Error() string {
	return "Route Not Found"
}
