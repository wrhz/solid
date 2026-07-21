package wasm

import "syscall/js"

func getGoExports() js.Value {
	goExports := js.Global().Get("goExports")

	if goExports.IsUndefined() || goExports.IsNull() {
		js.Global().Set("goExports", js.Global().Get("Map").New())
		goExports = js.Global().Get("goExports")
	}

	return goExports
}