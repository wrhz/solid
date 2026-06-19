package route

import (
	"fmt"
	"solid/solid"
)

type Hello struct{}

func (h *Hello) RegisterRoute(r *solid.RouteStruct) {
	r.Get("/hello", h.helloGet)
	r.Websocket("/ws", h.wsWebSocket)
}

func (h *Hello) Init(r *solid.RouteStruct) {
	
}	

func (h *Hello) RegisterMiddleware(r *solid.RouteStruct) {
	
}

func (h *Hello) helloGet(c *solid.Context) {
	solid.ViewResponse(c, "index.html", 200)
}

func (h *Hello) wsWebSocket(websocket *solid.WebSocket) {
	for !websocket.IsClosed() {
		_, message, err := websocket.ReadMessage()

		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received message: %s\n", message)

		err = websocket.SendTextMessage("Hello from WebSocket!")

		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func NewHello() *Hello {
	return &Hello{}
}