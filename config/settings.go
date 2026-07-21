package config

import "github.com/gorilla/sessions"

type SettingsConfigStruct struct {
	staticMaxAge int

	maxBytesMemory         int64
	multipartFormMaxMemory int64

	sessionsPairs []byte
	sessionStore  sessions.Store

	trustedProxies []string
}

func NewSettingsConfig() *SettingsConfigStruct {
	return &SettingsConfigStruct{
		maxBytesMemory:         64 << 20,
		multipartFormMaxMemory: 32 << 20,
	}
}

func (s *SettingsConfigStruct) SetStaticMaxAge(staticMaxAge int) {
	s.staticMaxAge = staticMaxAge
}

func (s *SettingsConfigStruct) SetMaxBytesMemory(maxBytesMemory int64) {
	s.maxBytesMemory = maxBytesMemory
}

func (s *SettingsConfigStruct) SetMultipartFormMaxMemory(maxMemory int64) {
	s.multipartFormMaxMemory = maxMemory
}

func (s *SettingsConfigStruct) SetSessionsSecret(sessionsPairs ...string) {
	for _, pair := range sessionsPairs {
		s.sessionsPairs = append(s.sessionsPairs, []byte(pair)...)
	}
}

func (s *SettingsConfigStruct) SetSessionStore(sessionStore sessions.Store) {
	s.sessionStore = sessionStore
}

func (s *SettingsConfigStruct) SetTrustedProxies(trustedProxies []string) {
	s.trustedProxies = trustedProxies
}

func (s *SettingsConfigStruct) GetStaticMaxAge() int {
	return s.staticMaxAge
}

func (s *SettingsConfigStruct) GetMaxBytesMemory() int64 {
	return s.maxBytesMemory
}

func (s *SettingsConfigStruct) GetMultipartFormMaxMemory() int64 {
	return s.multipartFormMaxMemory
}

func (s *SettingsConfigStruct) GetSessionsSecret() []byte {
	if s.sessionsPairs == nil {
		s.sessionsPairs = []byte{}
	}

	return s.sessionsPairs
}

func (s *SettingsConfigStruct) GetSessionStore() sessions.Store {
	return s.sessionStore
}

func (s *SettingsConfigStruct) GetTrustedProxies() []string {
	if s.trustedProxies == nil {
		s.trustedProxies = []string{}
	}

	return s.trustedProxies
}
