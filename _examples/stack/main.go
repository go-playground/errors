package main

import (
	"fmt"

	"strings"

	"github.com/go-playground/errors"
)

func main() {
	// maybe you just want to grab a stack trace and process on your own like go-playground/log
	// uses it to produce a stack trace log message
	frame := errors.Stack()
	name := fmt.Sprintf("%n", frame)
	file := fmt.Sprintf("%+s", frame)
	line := fmt.Sprintf("%d", frame)
	parts := strings.Split(file, "\n\t")
	if len(parts) > 1 {
		file = parts[1]
	}

	fmt.Printf("Name: %s File: %s Line: %s\n", name, file, line)
}
