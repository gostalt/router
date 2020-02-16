package router

type Group struct {
	prefix string
	routes []*Route
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
