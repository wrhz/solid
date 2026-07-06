package config

import (
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

type SettingsConfigStruct struct {
	staticMaxAge int

	maxBytesMemory         int64
	multipartFormMaxMemory int64

	sessionsPairs []byte
	sessionStore  sessions.Store
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

func (s *SettingsConfigStruct) GetStaticMaxAge() (int, error) {
	staticMaxAge := s.staticMaxAge
	if staticMaxAge == 0 {
		value, exists := os.LookupEnv("SOLID_STATIC_MAX_AGE")
		if exists {
			if val, err := strconv.Atoi(value); err == nil {
				staticMaxAge = val
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	}

	return staticMaxAge, nil
}

func (s *SettingsConfigStruct) GetMaxBytesMemory() (int64, error) {
	maxBytesMemory := s.maxBytesMemory
	if maxBytesMemory == 0 {
		value, exists := os.LookupEnv("SOLID_MAX_BYTES_MEMORY")
		if exists {
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				maxBytesMemory = val
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	}

	return maxBytesMemory, nil
}

func (s *SettingsConfigStruct) GetMultipartFormMaxMemory() (int64, error) {
	multipartFormMaxMemory := s.multipartFormMaxMemory
	if multipartFormMaxMemory == 0 {
		value, exists := os.LookupEnv("SOLID_MULTIPART_FORM_MAX_MEMORY")
		if exists {
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				multipartFormMaxMemory = val
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	}

	return multipartFormMaxMemory, nil
}

func (s *SettingsConfigStruct) GetSessionStore() sessions.Store {
	return s.sessionStore
}