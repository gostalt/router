package router

type Group struct {
	prefix string
	routes []*Route

	middleware []Middleware
}

func NewGroup(routes ...*Route) *Group {
	return &Group{
		routes: routes,
	}
}

func (g *Group) Prefix(path string) *Group {
	g.prefix = path
	return g
}

func (g *Group) Middleware(middleware ...Middleware) *Group {
	g.middleware = middleware
	return g
}

func (g *Group) Add(routes ...*Route) *Group {
	g.routes = append(g.routes, routes...)
	return g
}

func (g *Group) Routes() []*Route {
	return g.routes
}
