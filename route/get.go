package method

import (
	"net/http"

	"github.com/wrhz/solid/server"
)

var getRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return getRoutes
}

func (r *RouteStruct) Get(path string, callFunc func(c *server.Context) error, middlewares ...func(c *server.Context)) {
    handler := routeFuncHandle(callFunc)

    for i := len(middlewares) - 1; i >= 0; i-- {
        currentMw := middlewares[i]
        currentHandler := handler

        handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            c := server.NewContext(req, w)

            c.SetNext(currentHandler.ServeHTTP)

            currentMw(c)
        })
    }

    getRoutes[r.perfix+path] = r.routeChain(handler).ServeHTTP
}