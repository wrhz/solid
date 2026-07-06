package server

import (
	"fmt"

	"github.com/goccy/go-reflect"
)

func (c *Context) SaveCookie(s any, option *CookieOption) error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("SaveCookie: expected struct, got %v", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		cookieTag := field.Tag.Get("cookie")

		if cookieTag != "" {
			value := v.Field(i)

			if isNil(value) {
				c.RemoveCookie(cookieTag)
			} else {
				c.SetCookie(&Cookie{
					Name:   cookieTag,
					Value:  fmt.Sprintf("%v", value.Interface()),
				}, option)
			}
		}
	}

	return nil
}

func (c *Context) SaveSession(s any, name string, options *SessionOptions) error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("SaveSession: expected struct, got %v", v.Kind())
	}

	session, err := c.Session(name, options)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		sessionTag := field.Tag.Get("session")
		
		if sessionTag != "" {
			value := v.Field(i)

			if isNil(value) {
				if err := session.RemoveValue(sessionTag); err != nil {
					return fmt.Errorf("failed to remove session %q: %w", sessionTag, err)
				}
			} else {
				if err := session.Set(sessionTag, value.Interface()); err != nil {
					return fmt.Errorf("failed to set session %q: %w", sessionTag, err)
				}
			}
		}
	}

	return nil
}