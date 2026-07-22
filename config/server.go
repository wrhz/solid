package config

import (
	"solid/route"

    "github.com/wrhz/solid/config"
)

func ServerConfig() {
    server := config.GetServerConfig()

    server.SetPort(8000)

    server.SetMainStruct(route.NewHello())
}