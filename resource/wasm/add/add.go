package main

import (
	"fmt"
	"syscall/js"

	"github.com/wrhz/solid/wasm"
)

var message = "Hello World"

type Struct struct {
	Name string `wasm:"name"`
	Age int `wasm:"age"`
}

func (s Struct) SayHello(this js.Value, args []js.Value) {
	fmt.Println("Hello Go Wasm")
}

func (s *Struct) SayInfo(this js.Value, args []js.Value) {
	fmt.Printf("Name: %s Age: %d\n", s.Name, s.Age)
}

func add(this js.Value, args []js.Value) any {
	a := args[0].Int()
	b := args[1].Int()

	return a + b
}

func main() {
	wasm.ExportFunc("add", add)

	wasm.ExportVar("message", message)

	wasm.ExportClass("Struct", Struct{ Name: "Tom", Age: 13 })

	select {}
}