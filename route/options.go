package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var optionsRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func OptionsRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return optionsRoutes
}

func (r *RouteStruct) Options(path string, callFunc func(c *server.Context) error, middlewares ...func(c *server.Context, next http.HandlerFunc)) {
    handler := routeFuncHandle(callFunc)

    for i := len(middlewares) - 1; i >= 0; i-- {
        currentMw := middlewares[i]
        currentHandler := handler

        handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            c := &server.Context{Writer: w, Request: req}
            currentMw(c, currentHandler.ServeHTTP)
        })
    }

    optionsRoutes[r.perfix+path] = r.routeChain(handler).ServeHTTP
}