package wasm

import "syscall/js"

func Object() js.Value {
	return js.Global().Get("Object")
}

func Console() js.Value {
	return js.Global().Get("console")
}