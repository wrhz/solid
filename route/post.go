package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var postRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func PostRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return postRoutes
}

func (r *RouteStruct) Post(path string, callFunc func(c *server.Context) error) {
	postRoutes[r.perfix+path] = r.routeChain(routeFuncHandle(callFunc)).ServeHTTP
}