package config

import (
	"database/sql"

	"gorm.io/gorm"
)

type IDatabaseConfig interface {
	GetGormDialector() gorm.Dialector
	GetGormModels() []any
	GetGormOptions() []gorm.Option
	GetXormDataSourceName() string
	GetXormDriverName() string
	GetXormDriverOptions() []func(db *sql.DB) error
	GetXormModels() []any
	GetXormShowSQL() bool
	RegisterGormModels(models ...any)
	RegisterXormModels(models ...any)
}