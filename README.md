# Router

Router is a standard library compatible package that simplifies creating routes
for web applications with features such as route parameters, groups, and
callbacks.

## The Basics

A new router is created with the `router.New()` function. The returned value is
compatible with the standard library's `ListenAndServe` call:

```go
r := router.New()

if err := http.ListenAndServe(":8080", r); err != nil {
    panic(err)
}
```

The most basic route definition is accepts a URI and a function that returns a
string. The returned string is raw HTML.

```go
r.Get("welcome", func() string {
    return "<h1>Hello World</h1>"
})
```

### Methods

The router allows you to register routes for any HTTP verb:

```go
r.Get(uri, callback)
r.Post(uri, callback)
r.Put(uri, callback)
r.Patch(uri, callback)
r.Delete(uri, callback)
```

If you need to register a route that responds to multiple verbs, you can use the
`Match` function on the router instance. If a route should respond to any verb,
you can use the `Any` function:

```go
r.Match([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, "endpoint", func() string {
    // ...
})

r.Any("endpoint", func() string {
    // ...
})
```

### Middleware

Chain a call to the `Middleware` function onto a route definition to wrap the
handler in a middleware. Note that middleware must match the following
signature:

```go
func(http.Handler) http.Handler
```

Middleware registered against the router executes first, followed by middleware
on the group and finally middleware on the specific route definition.

In the following code snippet, the middleware would be executed `one`, `two` and
finally `three`, before calling the route's handler:

```go
rtr := router.New()
rtr.Middleware(one)

rtr.Group(
    router.Get("/", handler).Middleware(three)
).Middleware(two)
```

### Redirect Routes

To define a route that redirects to another URI, you can use the `Redirect`
function on a router instance:

```go
r.Redirect("old", "new")
```

## Route Parameters

Sometimes, you may want to use a portion of the URL within your route — for
example, capturing a resource ID. You can do so by defining a route using
parameters. To define a route parameter, the name of the parameter should be
wrapped with `{}` curly braces.

Multiple parameters are supported on a single route, but the parameter names
should be unique:

```go
r.Get("posts/{postId}/comments/{commentId}", func(req *http.Request) string {
    return "Post " + req.Form.Get("postId")
})
```

The value of the parameter is injected into the requests `Form` variables, and
can be retrieved using `Form.Get`.

## Groups

Groups enable middleware and prefixes to be shared across a collection of
groups.

To create a new group, use the `Group` function and pass in a variadic list of
routes:

```go
r := router.New()

r.Group(
    router.Get("profile", handler),
    router.Get("comments", handler),
)
```

### Adding Middleware

To add middleware to all routes in a group, chain a call to the `Middleware`
function to the `Groups` call:

```go
r.Group(...).Middleware(ThrotteRequests)
```

### Adding Prefixes

To add a route prefix to all routes in a group, chain a call to the `Prefix`
function:

```go
r.Group(
    router.Get("users", func() string {
        // Would match the URL `/admin/users`
    })
).Prefix("admin")
```

## Handler Shapes

By default, the below handler shapes are supported, meaning that they can be
used when registering a new route.

```go
func() string

func(*http.Request) string

http.HandlerFunc

func(http.ResponseWriter, *http.Request)

http.Handler
```
