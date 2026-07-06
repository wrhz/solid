package config

import "github.com/wrhz/solid/types"

var configManager *ConfigManager

type ConfigManager struct {
	settingsConfig  types.ISettingsConfig
	websocketConfig types.IWebSocketConfig
	databaseConfig  types.IDatabaseConfig
}

func NewConfigManager(settingsConfig types.ISettingsConfig, websocketConfig types.IWebSocketConfig, databaseConfig types.IDatabaseConfig) {
	configManager = &ConfigManager{
		settingsConfig:  settingsConfig,
		websocketConfig: websocketConfig,
		databaseConfig:  databaseConfig,
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