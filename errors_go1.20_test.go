//go:build go1.20
// +build go1.20

package errors

import (
	"io"
	"testing"
)

func TestJoin(t *testing.T) {
	err1 := io.EOF
	err2 := io.ErrUnexpectedEOF

	err := Join(err1, err2)
	innerErr, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected Join to return an error that implements Unwrap() []error")
	}
	errs := innerErr.Unwrap()
	if len(errs) != 2 {
		t.Fatalf("expected Join to return an error that implements Unwrap() []error to return 2 errors")
	}
	if !Is(errs[0], io.EOF) {
		t.Fatalf("expected Join to return an error that implements Unwrap() []error to return io.EOF")
	}
	if !Is(errs[1], io.ErrUnexpectedEOF) {
		t.Fatalf("expected Join to return an error that implements Unwrap() []error to return io.ErrUnexpectedEOF")
	}

	// test wrapping and then unwrapping with Chain
	err = Wrap(err, "my test wrapped error")
	if !Is(err, io.EOF) {
		t.Fatalf("expected wrapped error to traverse into joined inner error EOF")
	}
	if !Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected wrapped error to traverse into joined inner error ErrUnexpectedEOF")
	}
}
