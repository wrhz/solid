package server

import (
	"github.com/wrhz/solid/server/binding"
	"google.golang.org/protobuf/proto"
)

func (c *Context) ShouldBindJSON(s any) error {
	return binding.JSONBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindXML(s any) error {
	return binding.XMLBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindProtobuf(s proto.Message) error {
	return binding.ProtobufBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindMsgPack(s any) error {
	return binding.MsgPackBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindYAML(s any) error {
	return binding.YAMLBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindTOML(s any) error {
	return binding.TOMLBinding{}.Bind(c.Request, s)
}

func (c *Context) MustBindJSON(s any) error  {
	if err := (binding.JSONBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindXML(s any) error {
	if err := (binding.XMLBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindProtobuf(s proto.Message) error {
	if err := (binding.ProtobufBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindMsgPack(s any) error  {
	if err := (binding.MsgPackBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindYAML(s any) error {
	if err := (binding.YAMLBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindTOML(s any) error {
	if err := (binding.TOMLBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}