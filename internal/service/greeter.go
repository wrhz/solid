package service

import (
    "context"
    "fmt"
)

type GreeterService interface {
    SayHello(ctx context.Context, name string) (string, error)
}

type greeterService struct{}

func NewGreeterService() GreeterService {
    return &greeterService{}
}

func (s *greeterService) SayHello(ctx context.Context, name string) (string, error) {
    return fmt.Sprintf("Hello %s", name), nil
}