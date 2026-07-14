package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var headRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func HeadRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return headRoutes
}

func (r *RouteStruct) Head(path string, callFunc func(c *server.Context) error, middlewares ...func(c *server.Context, next http.HandlerFunc)) {
    handler := routeFuncHandle(callFunc)

    for i := len(middlewares) - 1; i >= 0; i-- {
        currentMw := middlewares[i]
        currentHandler := handler

        handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            c := &server.Context{Writer: w, Request: req}
            currentMw(c, currentHandler.ServeHTTP)
        })
    }

    headRoutes[r.perfix+path] = r.routeChain(handler).ServeHTTP
}