package webtype

func (app WebApp) Reverse(routeName string, pairs ...string) string {
	route := app.Router.Get(routeName)
	if route == nil {
		panic("Can't find reverse route for '" + routeName + "'")
	}
	url, err := app.Router.Get(routeName).URL(pairs...)
	if err != nil {
		panic("Can't find reverse route matching parameters for '" + routeName + "'")
	}
	return url.String()
}
