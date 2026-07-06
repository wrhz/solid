package wasm

import (
	"os"
	"syscall/js"
)

var moduleName = os.Args[0]

func Export(name string, callFunc func(this js.Value, args []js.Value) any) {
	goExports := js.Global().Get("go_exports")

	if !goExports.Call("has", js.ValueOf(moduleName)).Bool() {
		goExports.Call("set", js.ValueOf(moduleName), js.Global().Get("Map").New())
	}

	goExports.Call("get", js.ValueOf(moduleName)).Call("set", js.ValueOf(name), js.FuncOf(callFunc))
}