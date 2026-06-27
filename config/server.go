package config

import (
	"solid/route"

    "github.com/wrhz/Solid"
)

func ServerConfig() {
    server := solid.GetServerConfig()

    server.SetPort(8000)

    server.SetMainStruct(route.NewHello())

    // server.SetTLSCertFile("cert.pem")
	// server.SetTLSKeyFile("key.pem")
}