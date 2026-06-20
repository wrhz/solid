package main

import (
	"fmt"
	"net/http"
	"solid/config"
	"solid/solid"
	"strconv"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConfigs() {
	config.ServerConfig()
	config.SettingsConfig()
	config.WebSocketConfig()
}

func handleRoutes(serve *mux.Router) {
	for path, callFunc := range solid.GetRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("GET")
	}

	for path, callFunc := range solid.PostRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("POST")
	}

	for path, callFunc := range solid.PatchRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("PATCH")
	}

	for path, callFunc := range solid.DeleteRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("DELETE")
	}

	for path, callFunc := range solid.PutRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("PUT")
	}

	for path, callFunc := range solid.OptionsRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("OPTIONS")
	}

	for path, callFunc := range solid.HeadRoutes() {
		serve.Handle(path, gziphandler.GzipHandler(http.HandlerFunc(callFunc))).Methods("HEAD")
	}

	for path, callFunc := range solid.WebsocketRoutes() {
		serve.Handle(path, http.HandlerFunc(callFunc))
	}
}

func main() {
	handleConfigs()

	serverConfig := solid.GetServerConfig()

	serve := mux.NewRouter()

	route := solid.NewRoute()

	serverConfig.GetMainStruct().Init(route)

	serverConfig.GetMainStruct().RegisterMiddleware(route)

	serverConfig.GetMainStruct().RegisterRoute(route)

	handleRoutes(serve)

	serve.PathPrefix("/static/ts/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/ts/", http.FileServer(http.Dir("./dist/static/ts"))),
		),
	)

	serve.PathPrefix("/static/js/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/js/", http.FileServer(http.Dir("./dist/static/js"))),
		),
	)

	serve.PathPrefix("/static/scss/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/scss/", http.FileServer(http.Dir("./dist/static/scss"))),
		),
	)

	serve.PathPrefix("/static/sass/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/sass/", http.FileServer(http.Dir("./dist/static/sass"))),
		),
	)

	serve.PathPrefix("/static/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
		),
	)

	serve.PathPrefix("/lib/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/lib/", http.FileServer(http.Dir("./dist/resource/lib"))),
		),
	)

	fmt.Println("Server starting on port:", serverConfig.GetPort())

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(serverConfig.GetPort()),
		Handler: serve,
	}

	if tlsConfig := solid.GetServerConfig().GetTLSConfig(); tlsConfig != nil {
		server.TLSConfig = tlsConfig
	}

	if certFile := solid.GetServerConfig().GetTLSCertFile(); certFile != "" {
		keyFile := solid.GetServerConfig().GetTLSKeyFile()

		if err := server.ListenAndServeTLS("./certs/" + certFile, "./certs/" + keyFile); err != nil {
			fmt.Println("Server failed:", err)
		}
		return
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server failed:", err)
	}
}
