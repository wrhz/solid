package main

import (
	"syscall/js"

	"github.com/wrhz/solid/wasm"
)

func add(_ js.Value, args []js.Value) any {
	a := args[0].Int()
	b := args[1].Int()
	return a + b
}

func main() {
	wasm.Export("add", add)

	select {}
}