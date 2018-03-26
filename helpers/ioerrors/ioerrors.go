package ioerrors

import (
	"io"

	"github.com/go-playground/errors"
)

// IOErrors helps classify io related errors
func IOErrors(w *errors.Wrapped, err error) (cont bool) {
	switch err {
	case io.EOF:
		w.WithTypes("io")
		return
	case io.ErrClosedPipe:
		w.WithTypes("Permanent", "io")
		return
	case io.ErrNoProgress:
		w.WithTypes("Permanent", "io")
		return
	case io.ErrShortBuffer:
		w.WithTypes("Permanent", "io")
		return
	case io.ErrShortWrite:
		w.WithTypes("Permanent", "io")
		return
	case io.ErrUnexpectedEOF:
		w.WithTypes("Transient", "io")
		return
	}
	return true
}
