package route

import "solid/solid"

type Hello struct{}

func (h *Hello) RegisterRoute(r *solid.RouteStruct) {
	r.Get("/hello", h.helloGet)
}

func (h *Hello) Init(r *solid.RouteStruct) {
	
}	

func (h *Hello) RegisterMiddleware(r *solid.RouteStruct) {
	
}

func (h *Hello) helloGet(c *solid.Context) error {
	return solid.HtmlViewResponse(c, "index", 200)
}

func NewHello() *Hello {
	return &Hello{}
}