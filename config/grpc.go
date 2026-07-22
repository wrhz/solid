package config

import "github.com/wrhz/solid/config"

func GrpcConfig() {
	grpcConfig := config.GetGrpcConfig()

	grpcConfig.UseGrpc(true)
}