package config

import (
	"github.com/wrhz/Solid"

	"github.com/gorilla/sessions"
)

func SettingsConfig() {
	settings := solid.GetSettingsConfig()

	settings.SetStaticMaxAge(3600)

	settings.SetMaxBytesMemory(1 << 20)

	secret := []byte("your-sessions-key-must-16|24|32-bytes-long!!!")

	settings.SetSessionStore(sessions.NewCookieStore(secret))
}