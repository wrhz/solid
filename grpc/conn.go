package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/wrhz/solid/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GetGrpcConn(host string, port int) (*ClientGrpc, error) {
	var creds grpc.DialOption

	grpcConfig := config.GetGrpcConfig()

	caCert := grpcConfig.GetClientCaCertFile()
	certFile := grpcConfig.GetClientTlsCertFile()
	keyFile := grpcConfig.GetClientTlsKeyFile()

	if caCert == "" {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else if certFile == "" && keyFile == "" {
		_creds, err := credentials.NewClientTLSFromFile("ca.crt", "")
		if err != nil {
			return nil, err
		}

		creds = grpc.WithTransportCredentials(_creds)
	} else {
		cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
		caCert, _ := os.ReadFile("ca.crt")
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			ServerName:   "localhost",
		}
		creds = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), creds)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &ClientGrpc{ Conn: conn, Context: ctx, Cancel: cancel }, err
}