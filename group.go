package router

type Group struct {
	prefix string
	routes []*Route

	middleware []Middleware
}

func (g *Group) calculateRouteRegexs() {
	for _, r := range g.routes {
		r.regex = r.calculateRouteRegex()
	}
}

func NewGroup(routes ...*Route) *Group {
	g := &Group{}

	for _, r := range routes {
		r.group = g
		g.calculateRouteRegexs()
	}

	g.routes = routes

	return g
}

func (g *Group) Prefix(path string) *Group {
	g.prefix = path
	g.calculateRouteRegexs()
	return g
}

func (g *Group) Middleware(middleware ...Middleware) *Group {
	g.middleware = middleware
	return g
}

func (g *Group) Add(routes ...*Route) *Group {
	for _, r := range routes {
		i, found := g.findExistingRoute(r)
		if found {
			g.routes[i] = r
		} else {
			g.routes = append(g.routes, r)
		}
	}
	return g
}

func (g *Group) findExistingRoute(route *Route) (int, bool) {
	for i, r := range g.routes {
		if r.path == route.path && methodsMatch(r, route) {
			return i, true
		}
	}

	return -1, false
}

func (g *Group) Routes() []*Route {
	return g.routes
}
