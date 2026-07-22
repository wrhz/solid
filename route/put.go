package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var putRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func PutRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return putRoutes
}

func (r *RouteStruct) Put(path string, callFunc func(c *server.Context) error, middlewares ...func(c *server.Context, next http.HandlerFunc)) *Route {
    handler := routeFuncHandle(callFunc)

    for i := len(middlewares) - 1; i >= 0; i-- {
        currentMw := middlewares[i]
        currentHandler := handler

        handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            c := &server.Context{Writer: w, Request: req}
            currentMw(c, currentHandler.ServeHTTP)
        })
    }

    perfix := r.perfix+path

    putRoutes[perfix] = r.routeChain(handler).ServeHTTP

    return &Route{ perfix: perfix }
}