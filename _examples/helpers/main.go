package main

import (
	"fmt"
	"net"

	"github.com/go-playground/errors"
	"github.com/go-playground/errors/helpers/neterrors"
)

func main() {
	errors.RegisterHelper(neterrors.NETErrors)
	_, err := net.ResolveIPAddr("tcp", "foo")
	if err != nil {
		err = errors.Wrap(err, "failed to perform operation")
	}

	// all that extra context, types and tags captured for free
	// there are more helpers and you can even create your own.
	fmt.Println(err)
}
