package method

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	solidManager "github.com/wrhz/solid/manager"
	solidWebsocket "github.com/wrhz/solid/server/websocket"
)

var websocketRoutes = map[string]func(w http.ResponseWriter, r *http.Request){}

func WebsocketRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return websocketRoutes
}

func (r *RouteStruct) Websocket(path string, callFunc func(websocket *solidWebsocket.WebSocket)) {
	websocketConfig := solidManager.GetWebSocketConfig()
	
	websocketRoutes[r.perfix+path] = r.websocketChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocketConfig.GetUpgrader().Upgrade(w, req, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}

		websocketStruct := solidWebsocket.NewWebSocket(conn, req)
		
		defer websocketStruct.Close()

		ticker := time.NewTicker(time.Duration(websocketConfig.GetPingInterval()) * time.Second)

		defer ticker.Stop()

		go func() {
			for range ticker.C {
				err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Duration(websocketConfig.GetPongWait()) * time.Second))
				if err != nil {
					fmt.Println("WebSocket ping error:", err)
					websocketStruct.Close()
					return
				}
			}
		}()

		callFunc(websocketStruct)
	})).ServeHTTP
}