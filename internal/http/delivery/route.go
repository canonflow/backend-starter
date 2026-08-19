package delivery

type IRoute interface {
	Setup()
}

type RouteGroup struct {
	Routes []IRoute
}

func (r *RouteGroup) Register(route IRoute) {
	r.Routes = append(r.Routes, route)
}

func (r *RouteGroup) Wire() {
	for _, r := range r.Routes {
		r.Setup()
	}
}
