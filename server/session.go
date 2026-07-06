package server

import (
	"net/http"

	"github.com/gorilla/sessions"

	solidManager "github.com/wrhz/solid/manager"
)

type Session struct {
	session *sessions.Session

	w http.ResponseWriter
	r *http.Request
}

type SessionOptions struct {
    Path     string
    Domain   string
    MaxAge   int
    Secure   bool
    HttpOnly bool
    SameSite http.SameSite
}

func (c *Context) Session(name string, options *SessionOptions) (*Session, error) {
	settingsConfig := solidManager.GetSettingsConfig()

	session, err := settingsConfig.GetSessionStore().Get(c.Request, name)
	if err != nil {
		return nil, err
	}

	session.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}

	return &Session{ session: session, w: c.Writer, r: c.Request }, nil
}

func (s *Session) Set(name any, value any) (error) {
	s.session.Values[name] = value
	return s.session.Save(s.r, s.w)
}

func (s *Session) Get(name any) any {
	return s.session.Values[name]
}

func (s *Session) Delete() (error) {
	s.session.Options.MaxAge = -1
	return s.session.Save(s.r, s.w)
}

func (s *Session) Clear() (error) {
	s.session.Values = make(map[any]any)
	return s.session.Save(s.r, s.w)
}

func (s *Session) RemoveValue(name any) (error) {
	delete(s.session.Values, name)
	return s.session.Save(s.r, s.w)
}