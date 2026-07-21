package server

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"

	solidRoute "github.com/wrhz/solid/route"
)

func InitRoutes(serve *mux.Router) {
	for path, callFunc := range solidRoute.GetRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("GET")
	}

	for path, callFunc := range solidRoute.PostRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("POST")
	}

	for path, callFunc := range solidRoute.PatchRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("PATCH")
	}

	for path, callFunc := range solidRoute.DeleteRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("DELETE")
	}

	for path, callFunc := range solidRoute.PutRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("PUT")
	}

	for path, callFunc := range solidRoute.OptionsRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("OPTIONS")
	}

	for path, callFunc := range solidRoute.HeadRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("HEAD")
	}

	for path, callFunc := range solidRoute.WebsocketRoutes() {
		serve.Handle(path, http.HandlerFunc(callFunc))
	}
}