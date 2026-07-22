package route

import (
	"fmt"

	"github.com/wrhz/solid"
	"github.com/wrhz/solid/grpc"
	"github.com/wrhz/solid/server"

	pb "solid/proto/helloworld"

	solidRoute "github.com/wrhz/solid/route"
)

type Hello struct {
	geeterClient *grpc.ClientGrpc
}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Init() error {
	return nil
}

func (h *Hello) RegisterRoute(r *solidRoute.RouteStruct) {
	r.Get("/hello", h.helloGet)
}

func (h *Hello) RegisterMiddleware(m *solidRoute.MiddlewareStruct) {
	
}

func (h *Hello) ServerStart() error {
	geeterClient, err := grpc.GetGrpcConn("localhost", 50051)
	
	if err != nil {
		return err
	}

	h.geeterClient = geeterClient

	return nil
}

func (h *Hello) ServerEnd() error {
	return h.geeterClient.Close()
}

func (h *Hello) helloGet(c *server.Context) error {
	client := pb.NewGreeterClient(h.geeterClient.Conn)

	r, err := client.SayHello(h.geeterClient.Context, &pb.HelloRequest{Name: "王壬浩泽"})

	if err != nil {
		return err
	}

	fmt.Println(r.GetMessage())

	return c.JSON(solid.H{
		"message": "ok",
	}, 200)
}
