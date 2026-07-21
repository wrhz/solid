package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-reflect"
	"github.com/wrhz/solid/util"
)

type Cookie struct {
	Name   string
	Value  string
}

type CookieOption struct {
	Quoted      bool
	Path        string
	Domain      string
	Expires     time.Time
	RawExpires  string
	MaxAge      int
	Secure      bool
	HttpOnly    bool
	SameSite    http.SameSite
	Partitioned bool
	Raw         string
	Unparsed    []string
}

func (c *Context) SetCookie(cookie *Cookie, option *CookieOption) {
	cookies := &http.Cookie{
		Name:   cookie.Name,
		Value:  cookie.Value,
		Quoted: option.Quoted,
		Path:   option.Path,
		Domain: option.Domain,
		Expires:  option.Expires,
		RawExpires: option.RawExpires,
		MaxAge:   option.MaxAge,
		Secure:   option.Secure,
		HttpOnly: option.HttpOnly,
		SameSite: option.SameSite,
		Partitioned: option.Partitioned,
		Raw:    option.Raw,
		Unparsed: option.Unparsed,
	}

	http.SetCookie(c.Writer, cookies)
}

func (c *Context) GetCookie(name string) (*Cookie, *CookieOption, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return nil, nil, err
	}

	cookies := &Cookie{
		Name:   cookie.Name,
		Value:  cookie.Value,
	}

	options := &CookieOption{
		Quoted: cookie.Quoted,
		Path:   cookie.Path,
		Domain: cookie.Domain,
		Expires:  cookie.Expires,
		RawExpires: cookie.RawExpires,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
		SameSite: cookie.SameSite,
		Partitioned: cookie.Partitioned,
		Raw:    cookie.Raw,
		Unparsed: cookie.Unparsed,
	}

	return cookies, options, nil
}

func (c *Context) RemoveCookie(name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})
}

func (c *Context) BindCookie(s any) error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("BindCookie: expected struct, got %v", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		cookieTag := field.Tag.Get("cookie")

		if cookieTag == "" {
			cookieTag = util.LowerFirst(field.Name)
		}

		cookie, err :=c.Request.Cookie(cookieTag)
		if err != nil {
			return fmt.Errorf("failed to get cookie %q: %w", cookieTag, err)
		}

		if cookie != nil {
			value, err := util.ParseType(cookie.Value, field.Type.Kind())
			if err != nil {
				return fmt.Errorf("failed to parse cookie %q: %w", cookieTag, err)
			}

			v.Field(i).Set(reflect.ValueOf(value))
		}
	}

	return nil
}