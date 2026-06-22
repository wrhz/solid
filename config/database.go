package config

import (
	"solid/solid"

	"gorm.io/driver/sqlite"
)

func DatabaseConfig() {
	databaseConfig := solid.GetDatabaseConfig()

	databaseConfig.SetGormDialector(sqlite.Open("example.db"))
}