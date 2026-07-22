package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type ClientGrpc struct {
	Conn    *grpc.ClientConn
	Context context.Context
	Cancel  context.CancelFunc
}

func (c *ClientGrpc) Close() error {
	err := c.Conn.Close()
	c.Cancel()

	return err
}