package server

import "github.com/wrhz/solid/database"

func MigrateModels() error {
	if err := database.MigrateModels(); err != nil {
		return err
	}

	if err := database.SyncModels(); err != nil {
		return err
	}

	return nil
}