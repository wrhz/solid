package route

import (
	solidRoute "github.com/wrhz/solid/route"
	"github.com/wrhz/solid/server"
)

type Hello struct{}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Init(r *solidRoute.RouteStruct) {
	
}

func (h *Hello) RegisterRoute(r *solidRoute.RouteStruct) {
	r.Get("/hello", h.helloGet)
}

func (h *Hello) RegisterMiddleware(r *solidRoute.RouteStruct) {
}

func (h *Hello) ServerStart() {

}

func (h *Hello) ServerEnd() {

}

func (h *Hello) helloGet(c *server.Context) error {
	return c.HtmlViewResponse("index", "index.html", 200, map[string]any{
		"Name": "Tom",
		"Age": 13,
	})
}
