package config

import "github.com/wrhz/solid/types/config"

var configManager *ConfigManager

type ConfigManager struct {
	settingsConfig  config.ISettingsConfig
	websocketConfig config.IWebSocketConfig
	databaseConfig  config.IDatabaseConfig
	templateConfig  config.ITemplateConfig
	validatorConfig config.IValidatorConfig
	corsConfig config.ICorsConfig
}

func NewConfigManager(settingsConfig config.ISettingsConfig,
	websocketConfig config.IWebSocketConfig,
	databaseConfig config.IDatabaseConfig,
	templateConfig config.ITemplateConfig,
	validatorConfig config.IValidatorConfig,
	corsConfig config.ICorsConfig,
) {
	configManager = &ConfigManager{
		settingsConfig:  settingsConfig,
		websocketConfig: websocketConfig,
		databaseConfig:  databaseConfig,
		templateConfig: templateConfig,
		validatorConfig: validatorConfig,
		corsConfig: corsConfig,
	}
}

func GetSettingsConfig() config.ISettingsConfig {
	return configManager.settingsConfig
}

func GetWebSocketConfig() config.IWebSocketConfig {
	return configManager.websocketConfig
}

func GetDatabaseConfig() config.IDatabaseConfig {
	return configManager.databaseConfig
}

func GetTemplateConfig() config.ITemplateConfig {
	return configManager.templateConfig
}

func GetValidatorConfig() config.IValidatorConfig {
	return configManager.validatorConfig
}

func GetCorsConfig() config.ICorsConfig {
	return configManager.corsConfig
}