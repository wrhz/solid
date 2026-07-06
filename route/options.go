package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var optionsRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func OptionsRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return optionsRoutes
}

func (r *RouteStruct) Options(path string, callFunc func(c *server.Context) error) {
	optionsRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}