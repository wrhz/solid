package config

import solidManager "github.com/wrhz/solid/manager"

var serverConfig = NewServerConfig()
var settingsConfig = NewSettingsConfig()
var websocketConfig = NewWebSocketConfig()
var databaseConfig = NewDatabaseConfig()
var templateConfig = NewTemplateConfigStruct()

func GetServerConfig() *ServerConfigStruct {
	return serverConfig
}

func GetSettingsConfig() *SettingsConfigStruct {
	return settingsConfig
}

func GetWebSocketConfig() *WebSocketConfigStruct {
	return websocketConfig
}

func GetDatabaseConfig() *DatabaseConfigStruct {
	return databaseConfig
}

func GetTemplateConfig() *TemplateConfigStruct {
	return templateConfig
}

func InitConfigManager() {
	solidManager.NewConfigManager(
		settingsConfig,
		websocketConfig,
		databaseConfig,
		templateConfig,
	)
}