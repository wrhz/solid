package grpc

import (
	"google.golang.org/grpc"

	"solid/internal/grpc/handler"
	"solid/internal/service"

	pb "solid/proto/helloworld"
)

func InitServer(server *grpc.Server) {
	geeterService := service.NewGreeterService()

    pb.RegisterGreeterServer(server, handler.NewGreeterHandler(geeterService))
}