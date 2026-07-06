package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	*websocket.Conn

	Request *http.Request

	isClosed bool
}

func NewWebSocket(conn *websocket.Conn, req *http.Request) *WebSocket {
	return &WebSocket{
		Conn:     conn,
		Request:  req,
		isClosed: false,
	}
}

func (ws *WebSocket) Close() error {
	ws.isClosed = true
	return ws.Conn.Close()
}

func (ws *WebSocket) IsClosed() bool {
	return ws.isClosed
}

func (ws *WebSocket) SendTextMessage(message string) error {
	return ws.WriteJSON(map[string]string{ "type": "text", "requestId": ws.Request.Context().Value("requestId").(string), "data": message })
}

func (ws *WebSocket) SendBinaryMessage(message []byte) error {
	return ws.WriteJSON(map[string]any{ "type": "binary", "requestId": ws.Request.Context().Value("requestId").(string), "data": message })
}

func (ws *WebSocket) SendJSONMessage(v any) error {
	return ws.WriteJSON(map[string]any{ "type": "json", "requestId": ws.Request.Context().Value("requestId").(string), "data": v })
}
