package method

import (
	"net/http"

	"github.com/wrhz/solid/middleware"
)

func (r *RouteStruct) routeChain(handler http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	return http.HandlerFunc(middleware.RouteMiddleware(handler))
}

func (r *RouteStruct) websocketChain(handler http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	return http.HandlerFunc(middleware.WebSocketMiddleware(handler))
}