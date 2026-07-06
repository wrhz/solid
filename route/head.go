package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var headRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func HeadRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return headRoutes
}

func (r *RouteStruct) Head(path string, callFunc func(c *server.Context) error) {
	headRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}