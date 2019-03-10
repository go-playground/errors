package main

import (
	"fmt"
	"net"

	"github.com/go-playground/errors"
)

func main() {
	errors.RegisterHelper(MyCustomErrHandler)
	_, err := net.ResolveIPAddr("tcp", "foo")
	if err != nil {
		err = errors.Wrap(err, "failed to perform operation")
	}

	// all that extra context, types and tags captured for free
	// there are more helpers and you can even create your own.
	fmt.Println(err)
}

// MyCustomErrHandler helps classify my errors
func MyCustomErrHandler(c errors.Chain, err error) (cont bool) {
	switch err.(type) {
	case net.UnknownNetworkError:
		_ = c.AddTypes("io").AddTag("additional", "details")
		return
		//case net.Other:
		//	...
		//	return
	}
	return true
}
