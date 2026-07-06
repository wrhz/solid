package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var deleteRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func DeleteRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return deleteRoutes
}

func (r *RouteStruct) Delete(path string, callFunc func(c *server.Context) error) {
	deleteRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}