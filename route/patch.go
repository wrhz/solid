package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var patchRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func PatchRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return patchRoutes
}

func (r *RouteStruct) Patch(path string, callFunc func(c *server.Context) error) {
	patchRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}