package config

import solidManager "github.com/wrhz/solid/manager"

var serverConfig = NewServerConfig()
var settingsConfig = NewSettingsConfig()
var websocketConfig = NewWebSocketConfig()
var databaseConfig = NewDatabaseConfig()
var templateConfig = NewTemplateConfigStruct()
var validatorConfig = NewValidatorConfigStruct()
var corsConfig = NewCorsConfig()
var grpcConfig = NewGrpcConfig()

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

func GetValidatorConfig() *ValidatorConfigStruct {
	return validatorConfig
}

func GetCorsConfig() *CorsConfig {
	return corsConfig
}

func GetGrpcConfig() *GrpcConfig {
	return grpcConfig
}

func InitConfigManager() {
	solidManager.NewConfigManager(
		settingsConfig,
		websocketConfig,
		databaseConfig,
		templateConfig,
		validatorConfig,
		corsConfig,
	)
}