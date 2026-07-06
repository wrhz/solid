package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Upgrader struct {
	websocket.Upgrader

	HandshakeTimeout                time.Duration
	ReadBufferSize, WriteBufferSize int
	WriteBufferPool                 websocket.BufferPool
	Subprotocols                    []string
	Error                           func(w http.ResponseWriter, r *http.Request, status int, reason error)
	CheckOrigin                     func(r *http.Request) bool
	EnableCompression               bool

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
