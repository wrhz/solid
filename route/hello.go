package route

import (
	"github.com/wrhz/solid"
	"github.com/wrhz/solid/server"

	solidRoute "github.com/wrhz/solid/route"
)

type MessageStruct struct {
	Message string `json:"message"`
}

type Hello struct {}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Init(r *solidRoute.RouteStruct) {
	
}

func (h *Hello) RegisterRoute(r *solidRoute.RouteStruct) {
	r.Get("/hello", h.helloGet)
}

func (h *Hello) RegisterMiddleware(m *solidRoute.MiddlewareStruct) {
	
}

func (h *Hello) ServerStart() {

}

func (h *Hello) ServerEnd() {

}

func (h *Hello) helloGet(c *server.Context) error {
	return c.JSON(solid.H{
		"message": "ok",
	}, 200)
}
