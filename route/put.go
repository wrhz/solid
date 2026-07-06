package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var putRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func PutRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return putRoutes
}

func (r *RouteStruct) Put(path string, callFunc func(c *server.Context) error) {
	putRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}