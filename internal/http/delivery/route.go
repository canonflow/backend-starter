package delivery

type IRoute interface {
	Setup()
}

type RouteGroup struct {
	Routes []IRoute
}

func (r *RouteGroup) Register(routes ...IRoute) {
	r.Routes = routes
}

func (r *RouteGroup) Wire() {
	for _, r := range r.Routes {
		r.Setup()
	}
}
