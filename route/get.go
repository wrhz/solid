package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var getRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return getRoutes
}

func (r *RouteStruct) Get(path string, callFunc func(c *server.Context) error) {
	getRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}