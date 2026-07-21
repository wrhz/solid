package server

import (
	"fmt"

	solidManager "github.com/wrhz/solid/manager"
	"github.com/wrhz/solid/server/binding"
)

func (c *Context) ShouldBindQuery(s any) error {
	return binding.QueryBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindParam(s any) error {
	return binding.ParamBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindForm(s any) error {
	return binding.FormBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBindHeader(s any) error {
	return binding.HeaderBinding{}.Bind(c.Request, s)
}

func (c *Context) ShouldBind(s any) error {
	contentType := c.Request.Header["Content-Type"][0]
	binding, ok := binding.Bindings[contentType]

	if ok {
		return binding.Bind(c.Request, s)
	}

	return fmt.Errorf("We can't find Content-Type: " + contentType)
}

func (c *Context) ShouldBindWith(s any, b binding.Binding) error {
	return b.Bind(c.Request, s)
}

func (c *Context) ShouldBindBodyWith(s any, b binding.BindingBody) error {
	body, err := c.Body()

	if err != nil {
		return err
	}

	return b.BindBody(body, s)
}

func (c *Context) MustBindQuery(s any) error {
	if err := (binding.QueryBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindParam(s any) error {
	if err := (binding.ParamBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindForm(s any) error {
	if err := (binding.FormBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBindHeader(s any) error {
	if err := (binding.HeaderBinding{}.Bind(c.Request, s)); err != nil {
		c.Abort()

		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) MustBind(s any) error {
	contentType := c.Request.Header["Content-Type"][0]
	binding, ok := binding.Bindings[contentType]

	if ok {
		c.Abort()

		err := binding.Bind(c.Request, s)

		if err != nil {
			c.Abort()

			return c.JSONError(400, err)
		}

		return nil
	}

	c.Abort()

	return c.JSONError(400, fmt.Errorf("We can't find Content-Type: %s", contentType))
}

func (c *Context) MustBindWith(s any, b binding.Binding) error {
	err := b.Bind(c.Request, s)

	if err != nil {
		return c.JSONError(400, err)
	}

	return nil
}

func (c *Context) Validate(s any) error {
	if err := solidManager.GetValidatorConfig().GetValidator().Struct(s); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	
	return nil
}