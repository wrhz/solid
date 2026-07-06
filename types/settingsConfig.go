package types

import "github.com/gorilla/sessions"

type ISettingsConfig interface {
	GetMaxBytesMemory() (int64, error)
	GetMultipartFormMaxMemory() (int64, error)
	GetSessionStore() sessions.Store
	GetStaticMaxAge() (int, error)
}