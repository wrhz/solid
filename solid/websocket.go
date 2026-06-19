package solid

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	*websocket.Conn

	Request *http.Request

	isClosed bool
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

type Upgrader struct {
	websocket.Upgrader

	HandshakeTimeout time.Duration
	ReadBufferSize, WriteBufferSize int
	WriteBufferPool websocket.BufferPool
	Subprotocols []string
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)
	CheckOrigin func(r *http.Request) bool
	EnableCompression bool

	isFilled bool
}

func (u *Upgrader) fill() {
	u.isFilled = true
	u.Upgrader.HandshakeTimeout = u.HandshakeTimeout
	u.Upgrader.ReadBufferSize = u.ReadBufferSize
	u.Upgrader.WriteBufferSize = u.WriteBufferSize
	u.Upgrader.WriteBufferPool = u.WriteBufferPool
	u.Upgrader.Subprotocols = u.Subprotocols
	u.Upgrader.Error = u.Error
	u.Upgrader.CheckOrigin = u.CheckOrigin
	u.Upgrader.EnableCompression = u.EnableCompression
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	if !u.isFilled {
		u.fill()
	}
	conn, err := u.Upgrader.Upgrade(w, r, responseHeader)
	return conn, err
}
