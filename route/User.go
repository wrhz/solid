package route

import "github.com/wrhz/Solid"

type User struct {
	// Write your members
}

func NewUser() *User {
	// Write your code to new this struct

	return &User{}
}

func (user *User) Init(r *solid.RouteStruct) {
	// Write your init code
}

func (user *User) RegisterRoute(r *solid.RouteStruct) {
	// Register your routes
}	

func (user *User) RegisterMiddleware(r *solid.RouteStruct) {
	// Register your moddlewares
}

func (user *User) ServerStart() {
	// Write your code, and it will run when the server starts.
}

func (user *User) ServerEnd() {
		Write your code, and it will run when the server finishes.
}