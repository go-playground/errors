package main

import (
	"fmt"

	runtimeext "github.com/go-playground/pkg/v5/runtime"
)

func main() {
	// maybe you just want to grab a stack trace and process on your own like go-playground/log
	// uses it to produce a stack trace log message
	frame := runtimeext.Stack()
	fmt.Printf("Function: %s File: %s Line: %d\n", frame.Function(), frame.File(), frame.Line())

	// and still have access to the underlying runtime.Frame
	fmt.Printf("%+v\n", frame.Frame)
}
