//go:build go1.20
// +build go1.20

package errors

import stderrors "errors"

// Join allows this library to be a drop-in replacement to the std library.
//
// Join returns an error that wraps the given errors. Any nil error values are discarded.
// Join returns nil if every value in errs is nil. The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
//
// It is the responsibility of the caller to then check for nil and wrap this error if desired.
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}
