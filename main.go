package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"solid/config"
	"strconv"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/wrhz/solid/database"

	solidConfig "github.com/wrhz/solid/config"
	solidRoute "github.com/wrhz/solid/route"
)

var mainStruct solidRoute.SolidMainRoute

func parseFlags() {
	serverConfig := solidConfig.GetServerConfig()
	
	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	serverConfig.SetDebug(*debug)
}

func handleConfigs() {
	config.ServerConfig()
	config.SettingsConfig()
	config.WebSocketConfig()
	config.DatabaseConfig()

	solidConfig.InitConfigManager()
}

func handleRoutes(serve *mux.Router) {
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

func migrateModels() error {
	if err := database.MigrateModels(); err != nil {
		return err
	}

	if err := database.SyncModels(); err != nil {
		return err
	}

	return nil
}

func serverFinish() {
	mainStruct.ServerEnd()

	if err := database.RemoveGorm(); err != nil {
		fmt.Printf("Remove GORM error: %v\n", err)
	}

	if err := database.RemoveXorm(); err != nil {
		fmt.Printf("Remove XORM error: %v\n", err)
	}
}

func main() {
	parseFlags()

	handleConfigs()

	serverConfig := solidConfig.GetServerConfig()

	debug := serverConfig.GetDebug()

	serve := mux.NewRouter()

	route := solidRoute.NewRoute()

	middleware := solidRoute.NewMiddleware(route.GetMiddlewares())

	mainStruct = serverConfig.GetMainStruct()

	mainStruct.Init(route)

	mainStruct.RegisterMiddleware(middleware)

	mainStruct.RegisterRoute(route)

	handleRoutes(serve)

	handleStatic(serve)

	err := database.InitGorm()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = database.InitXorm()

	if err != nil {
		fmt.Println(err)
		return
	}

	if debug {
		if err = migrateModels(); err != nil {
			fmt.Println("Migrate Models error: ", err)
			return
		}
	}

	fmt.Println("Server starting on port:", serverConfig.GetPort())

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(serverConfig.GetPort()),
		Handler: serve,
	}

	if tlsConfig := solidConfig.GetServerConfig().GetTLSConfig(); tlsConfig != nil {
		httpServer.TLSConfig = tlsConfig
	}

	go func ()  {
		mainStruct.ServerStart()

		if certFile := solidConfig.GetServerConfig().GetTLSCertFile(); certFile != "" {
			keyFile := solidConfig.GetServerConfig().GetTLSKeyFile()

			if err := httpServer.ListenAndServeTLS("./certs/" + certFile, "./certs/" + keyFile); err != nil {
				fmt.Println("Server failed:", err)
			}
		} else {
			if err := httpServer.ListenAndServe(); err != nil {
				fmt.Println("Server failed:", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serverFinish()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server close timeout or error: %v\n", err)
	} else {
		fmt.Println("Server closed")
	}
}
