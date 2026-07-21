package server

import (
	"solid/config"

	solidConfig "github.com/wrhz/solid/config"
)

func InitConfigs() {
	config.ServerConfig()
	config.SettingsConfig()
	config.WebSocketConfig()
	config.DatabaseConfig()
	config.TemplateConfig()
	config.ValidatorConfig()
	config.CorsConfig()

	solidConfig.InitConfigManager()
}