package main

import (
	"fmt"
	"io"

	"github.com/go-playground/errors/v5"
	nestedpackagee "github.com/go-playground/errors/v5/_examples/basic/nestedpackage"
)

func main() {
	err := level1Function("1")
	fmt.Println(err)
	if errors.HasType(err, "Permanent") {
		// os.Exit(1)
		fmt.Println("it is a permanent error")
	}

	// root error
	cause := errors.Cause(err)
	fmt.Println("CAUSE:", cause)

	// can even still inspect the internal error
	fmt.Println(errors.Cause(err) == io.EOF) // will extract the cause for you
	fmt.Println(errors.Cause(cause) == io.EOF)

	// and still in a switch
	switch errors.Cause(err) {
	case io.EOF:
		fmt.Println("EOF error")
	default:
		fmt.Println("unknown error")
	}
}

func level1Function(userID string) error {
	if err := level2Function(userID); err != nil {
		return errors.Wrap(err, "level2 call failed")
	}
	return nil
}

func level2Function(userID string) error {
	err := nestedpackagee.GetUser(userID)
	return errors.Wrap(err, "GetUser something went wrong")
}
