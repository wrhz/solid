package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"solid/config"
	"solid/solid"
	"strconv"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
)

func handleConfigs() {
	config.ServerConfig()
	config.SettingsConfig()
	config.WebSocketConfig()
	config.DatabaseConfig()
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

func handleStatic(serve *mux.Router) {
	serve.PathPrefix("/static/").Handler(
		gziphandler.GzipHandler(
			http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
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

func serverFinish() {
	if err := solid.RemoveGorm(); err != nil {
		fmt.Printf("Remove GORM error: %v\n", err)
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

	handleStatic(serve)

	solid.InitGorm()

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

	go func ()  {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server close timeout or error: %v\n", err)
	} else {
		fmt.Println("Server closed")
	}

	serverFinish()
}
