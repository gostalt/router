package router

import "net/http"

// Middleware are handlers that are added as part of a route definition. They are
// used to wrap the route's handler in additional layers of logic.
type Middleware func(http.Handler) http.Handler
