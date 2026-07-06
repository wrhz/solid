package server

import (
	"net/http"
	"time"
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