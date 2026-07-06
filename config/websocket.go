package config

import (
	"net/http"

	"github.com/wrhz/solid/server/websocket"
)

type WebSocketConfigStruct struct {
	upgrader *websocket.Upgrader

	pingInterval int
	pongWait     int
}

func (w *WebSocketConfigStruct) SetUpgrader(upgrader *websocket.Upgrader) {
	w.upgrader = upgrader
}

func (w *WebSocketConfigStruct) SetPingInterval(pingInterval int) {
	w.pingInterval = pingInterval
}

func (w *WebSocketConfigStruct) SetPongWait(pongWait int) {
	w.pongWait = pongWait
}

func (w *WebSocketConfigStruct) GetUpgrader() *websocket.Upgrader {
	return w.upgrader
}

func (w *WebSocketConfigStruct) GetPingInterval() int {
	return w.pingInterval
}

func (w *WebSocketConfigStruct) GetPongWait() int {
	return w.pongWait
}

func NewWebSocketConfig() *WebSocketConfigStruct {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WebSocketConfigStruct{
		upgrader:     upgrader,
		pingInterval: 10,
		pongWait:     5,
	}
}