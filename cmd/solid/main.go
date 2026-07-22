package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/wrhz/solid/database"
	"golang.org/x/crypto/acme/autocert"

	solidConfig "github.com/wrhz/solid/config"
	solidInit "github.com/wrhz/solid/init"
)

func parseFlags() {
	serverConfig := solidConfig.GetServerConfig()
	
	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	serverConfig.SetDebug(*debug)
}

func main() {
	parseFlags()

	serverConfig := solidConfig.GetServerConfig()

	mainStruct, serve, err := solidInit.InitServer(serverConfig)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server starting on port:", serverConfig.GetPort())

	httpServer := serverConfig.GetServerConfig()

	tlsConfig := serverConfig.GetTLSConfig()
	if tlsConfig != nil {
		httpServer.TLSConfig = tlsConfig
	} else if serverConfig.GetAutoTLS() {
		m := &autocert.Manager{
			Cache: autocert.DirCache("./.cache"),

			HostPolicy: autocert.HostWhitelist(serverConfig.GetAutoTLSHostPolicy()...),
		}

		go func() {
			http.ListenAndServe(":80", m.HTTPHandler(nil))
		}()

		httpServer.TLSConfig = &tls.Config{
			GetCertificate: m.GetCertificate,
		}
	}

	httpServer.Addr = ":" + strconv.Itoa(serverConfig.GetPort())
	httpServer.Handler = serve

	mainStruct.ServerStart()

	go func ()  {
		if solidConfig.GetServerConfig().GetAutoTLS() {
			httpServer.ListenAndServeTLS("", "")
		} else if certFile := serverConfig.GetTLSCertFile(); certFile != "" {
			keyFile := serverConfig.GetTLSKeyFile()

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

	mainStruct.ServerEnd()

	if err := database.RemoveGorm(); err != nil {
		fmt.Printf("Remove GORM error: %v\n", err)
	}

	if err := database.RemoveXorm(); err != nil {
		fmt.Printf("Remove XORM error: %v\n", err)
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server close timeout or error: %v\n", err)
	} else {
		fmt.Println("Server closed")
	}
}
