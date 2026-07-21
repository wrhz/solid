package config

import "github.com/gorilla/sessions"

type ISettingsConfig interface {
	GetMaxBytesMemory() int64
	GetMultipartFormMaxMemory() int64
	GetSessionStore() sessions.Store
	GetSessionsSecret() []byte
	GetStaticMaxAge() int
	GetTrustedProxies() []string
}