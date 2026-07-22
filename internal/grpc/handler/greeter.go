package handler

import (
    "context"
    "solid/internal/service"
    pb "solid/proto/helloworld"
)

type GreeterHandler struct {
    pb.UnimplementedGreeterServer
    svc service.GreeterService
}

func NewGreeterHandler(svc service.GreeterService) *GreeterHandler {
    return &GreeterHandler{svc: svc}
}

func (h *GreeterHandler) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    msg, err := h.svc.SayHello(ctx, req.GetName())
    if err != nil {
        return nil, err
    }
    return &pb.HelloReply{Message: msg}, nil
}