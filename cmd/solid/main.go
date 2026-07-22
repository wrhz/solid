package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/wrhz/solid/database"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	solidConfig "github.com/wrhz/solid/config"
	solidInit "github.com/wrhz/solid/init"

	projectGrpc "solid/internal/grpc"
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

	go func ()  {
		fmt.Println("HTTP server starting on port:", serverConfig.GetPort())

		if serverConfig.GetAutoTLS() {
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

	grpcConfig := solidConfig.GetGrpcConfig()

	if grpcConfig.GetUseGrpc() {
		go func ()  {
			var server *grpc.Server

			if grpcConfig.GetServerTlsCertFile() != "" && grpcConfig.GetClientTlsKeyFile() != "" {
				creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
				if err != nil {
					fmt.Printf("failed to load credentials: %v\n", err)
				}

				server = grpc.NewServer(grpc.Creds(creds))
			} else  {
				server = grpc.NewServer()
			}

    		projectGrpc.InitServer(server)

			fmt.Println("gRPC server starting on port:", grpcConfig.GetPort())

			lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcConfig.GetHost(), grpcConfig.GetPort()))
			if err != nil {
				fmt.Printf("failed to listen: %v\n", err)
			}

			if err := server.Serve(lis); err != nil {
				fmt.Printf("failed to serve: %v\n", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := mainStruct.ServerEnd(); err != nil {
		fmt.Printf("Run MainRoute error: %v\n", err)
	}

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
