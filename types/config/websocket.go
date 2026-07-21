package config

import "github.com/wrhz/solid/server/websocket"

type IWebSocketConfig interface {
	GetPingInterval() int
	GetPongWait() int
	GetUpgrader() *websocket.Upgrader
}