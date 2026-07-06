package database

import (
	"sync"

	"gorm.io/gorm"

	solidManager "github.com/wrhz/solid/manager"
)

var gormDatabase *gorm.DB
var GormDatabasesManager *GormDatabasesManagerStruct = &GormDatabasesManagerStruct{}

type GormDatabasesManagerStruct struct {
	store sync.Map
}

func (g *GormDatabasesManagerStruct) Set(requestID string) {
	g.store.Store(requestID, gormDatabase.Begin())
}

func (g *GormDatabasesManagerStruct) Get(requestID string) (*gorm.DB, bool) {
	orm, ok := g.store.Load(requestID)

	if ok {
		return orm.(*gorm.DB), true
	}

	return nil, false
}

func (g *GormDatabasesManagerStruct) Delete(requestID string) {
	g.store.Delete(requestID)
}

func InitGorm() error {
	databaseConfig := solidManager.GetDatabaseConfig()

	gormDialector := databaseConfig.GetGormDialector()

	if gormDialector != nil {
		var err error

		gormOptions := databaseConfig.GetGormOptions()
		gormDatabase, err = gorm.Open(gormDialector, gormOptions...)

		return err
	}

	return nil
}

func RemoveGorm() error {
	if IsStartGorm() {
		sqlDB, err := gormDatabase.DB()

		if err != nil {
			sqlDB.Close()
		} else {
			return nil
		}
	}

	return nil
}

func MigrateModels() error {
	if IsStartGorm() {
		databaseConfig := solidManager.GetDatabaseConfig()

		return gormDatabase.AutoMigrate(databaseConfig.GetGormModels()...)
	}

	return nil
}

func IsStartGorm() bool {
	return gormDatabase != nil
}
