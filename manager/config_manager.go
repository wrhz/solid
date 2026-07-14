package config

import "github.com/wrhz/solid/types"

var configManager *ConfigManager

type ConfigManager struct {
	settingsConfig  types.ISettingsConfig
	websocketConfig types.IWebSocketConfig
	databaseConfig  types.IDatabaseConfig
	templateConfig  types.ITemplateConfig
}

func NewConfigManager(settingsConfig types.ISettingsConfig,
	websocketConfig types.IWebSocketConfig,
	databaseConfig types.IDatabaseConfig,
	templateConfig types.ITemplateConfig,
) {
	configManager = &ConfigManager{
		settingsConfig:  settingsConfig,
		websocketConfig: websocketConfig,
		databaseConfig:  databaseConfig,
		templateConfig: templateConfig,
	}
}

func GetSettingsConfig() types.ISettingsConfig {
	return configManager.settingsConfig
}

func GetWebSocketConfig() types.IWebSocketConfig {
	return configManager.websocketConfig
}

func GetDatabaseConfig() types.IDatabaseConfig {
	return configManager.databaseConfig
}

func GetTemplateConfig() types.ITemplateConfig {
	return configManager.templateConfig
}