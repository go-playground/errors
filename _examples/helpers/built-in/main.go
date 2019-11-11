package main

import (
	"fmt"
	"net"

	"github.com/go-playground/errors/v5"
	// init function handles registration automatically
	_ "github.com/go-playground/errors/v5/helpers/neterrors"
)

func main() {
	_, err := net.ResolveIPAddr("tcp", "foo")
	if err != nil {
		err = errors.Wrap(err, "failed to perform operation")
	}

	// all that extra context, types and tags captured for free
	// there are more helpers and you can even create your own.
	fmt.Println(err)
}
