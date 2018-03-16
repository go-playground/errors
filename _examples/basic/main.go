package main

import (
	"fmt"
	"io"

	"github.com/go-playground/errors"
)

func main() {
	err := level1("testing error")
	fmt.Println(err)
	if errors.HasType(err, "Permanent") {
		// os.Exit(1)
	}

	// root error
	cause := errors.Cause(err)
	fmt.Println(cause)

	// can even still inspect the internal error
	fmt.Println(errors.IsErr(err, io.EOF)) // will extract the cause for you
	fmt.Println(errors.IsErr(cause, io.EOF))

	// or manually with access to base fields
	wrapped := cause.(*errors.Wrapped)
	fmt.Println(wrapped.Err == io.EOF)
}

func level1(value string) error {
	if err := level2(value); err != nil {
		return errors.Wrap(err, "level2 call failed")
	}
	return nil
}

func level2(value string) error {
	err := fmt.Errorf("this is an %s", "error")
	return errors.Wrap(err, "failed to do something").WithTypes("Permanent").WithTags(errors.T("value", value))
}
