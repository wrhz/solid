package wasm

import (
	"sync"
	"syscall/js"
)

var initOnce sync.Once
var funcs = []js.Func{}

func addFunc(callFunc js.Func) js.Func {
	funcs = append(funcs, callFunc)

	return callFunc
}

func InitFunc() {
	initOnce.Do(func() {
		var destroyFunc js.Func

		destroyFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
			for _, callFunc := range funcs {
				callFunc.Release()
			}

			funcs = nil

			destroyFunc.Release()

			return nil
		})

		js.Global().Call("addEventListener", "unload", destroyFunc)
		js.Global().Call("addEventListener", "unload", destroyFunc)
	})
}