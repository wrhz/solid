package config

import (
	"database/sql"

	"gorm.io/gorm"
)

type DatabaseConfigStruct struct {
	gormDialector gorm.Dialector
	gormOptions   []gorm.Option
	gormModels    []any

	xormDriverName     string
	xormDataSourceName string
	xormDriverOptions  []func(db *sql.DB) error
	xormShowSQL        bool
	xormModels         []any
}

func (d *DatabaseConfigStruct) SetGormDialector(dialector gorm.Dialector) {
	d.gormDialector = dialector
}

func (d *DatabaseConfigStruct) SetGormOptions(options ...gorm.Option) {
	d.gormOptions = append(d.gormOptions, options...)
}

func (d *DatabaseConfigStruct) GetGormDialector() gorm.Dialector {
	return d.gormDialector
}

func (d *DatabaseConfigStruct) GetGormOptions() []gorm.Option {
	return d.gormOptions
}

func (d *DatabaseConfigStruct) RegisterGormModels(models ...any) {
	d.gormModels = append(d.gormModels, models...)
}

func (d *DatabaseConfigStruct) GetGormModels() []any {
	return d.gormModels
}

func (d *DatabaseConfigStruct) GetXormDriverName() string {
	return d.xormDriverName
}

func (d *DatabaseConfigStruct) SetXormDriverName(name string) {
	d.xormDriverName = name
}

func (d *DatabaseConfigStruct) GetXormDataSourceName() string {
	return d.xormDataSourceName
}

func (d *DatabaseConfigStruct) SetXormDataSourceName(name string) {
	d.xormDataSourceName = name
}

func (d *DatabaseConfigStruct) GetXormDriverOptions() []func(db *sql.DB) error {
	return d.xormDriverOptions
}

func (d *DatabaseConfigStruct) SetXormDriverOptions(opts []func(db *sql.DB) error) {
	d.xormDriverOptions = append(d.xormDriverOptions, opts...)
}

func (d *DatabaseConfigStruct) SetXormShowSQL(showSQL bool) {
	d.xormShowSQL = showSQL
}

func (d *DatabaseConfigStruct) GetXormShowSQL() bool {
	return d.xormShowSQL
}

func (d *DatabaseConfigStruct) RegisterXormModels(models ...any) {
	d.xormModels = append(d.xormModels, models...)
}

func (d *DatabaseConfigStruct) GetXormModels() []any {
	return d.xormModels
}

func NewDatabaseConfig() *DatabaseConfigStruct {
	return &DatabaseConfigStruct{}
}