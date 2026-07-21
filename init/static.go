package server

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
)

func InitStatic(serve *mux.Router) {
	serve.PathPrefix("/static/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
		),
	)

	serve.PathPrefix("/resource/wasm/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/resource/wasm/", http.FileServer(http.Dir("./output/wasm"))),
		),
	)

	serve.PathPrefix("/resource/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/resource/", http.FileServer(http.Dir("./dist/resource"))),
		),
	)

	serve.PathPrefix("/chunks/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/chunks/", http.FileServer(http.Dir("./dist/chunks"))),
		),
	)
}