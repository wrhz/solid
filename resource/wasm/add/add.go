package main

import (
	"fmt"
	"os"
	"syscall/js"
)

func add(this js.Value, args []js.Value) any {
	a := args[0].Int()
	b := args[1].Int()
	return a + b
}

func main() {
	fmt.Println(os.Args)

	select {}
}