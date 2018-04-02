package ioerrors

import (
	"io"

	"github.com/go-playground/errors"
)

// IOErrors helps classify io related errors
func IOErrors(c errors.Chain, err error) (cont bool) {
	switch err {
	case io.EOF:
		c.WithTypes("io")
		return
	case io.ErrClosedPipe:
		c.WithTypes("Permanent", "io")
		return
	case io.ErrNoProgress:
		c.WithTypes("Permanent", "io")
		return
	case io.ErrShortBuffer:
		c.WithTypes("Permanent", "io")
		return
	case io.ErrShortWrite:
		c.WithTypes("Permanent", "io")
		return
	case io.ErrUnexpectedEOF:
		c.WithTypes("Transient", "io")
		return
	}
	return true
}
