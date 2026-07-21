package server

import (
	"github.com/gorilla/mux"
	"github.com/wrhz/solid/database"
	"github.com/wrhz/solid/server"

	solidConfig "github.com/wrhz/solid/config"
	solidRoute "github.com/wrhz/solid/route"
)

func InitServer(serverConfig *solidConfig.ServerConfigStruct) (solidRoute.SolidMainRoute, *mux.Router, error) {
	InitConfigs()

	debug := serverConfig.GetDebug()

	serve := mux.NewRouter()

	route := solidRoute.NewRoute()

	middleware := solidRoute.NewMiddleware(route.GetMiddlewares())

	mainStruct := serverConfig.GetMainStruct()

	mainStruct.Init(route)

	mainStruct.RegisterMiddleware(middleware)

	mainStruct.RegisterRoute(route)

	InitRoutes(serve)

	InitStatic(serve)

	if err := server.LoadHTML(); err != nil {
		return nil, nil, err
	}

	if err := database.InitGorm(); err != nil {
		return nil, nil, err
	}

	if err := database.InitXorm(); err != nil {
		return nil, nil, err
	}

	if debug {
		if err := MigrateModels(); err != nil {
			return nil, nil, err
		}
	}

	return mainStruct, serve, nil
}