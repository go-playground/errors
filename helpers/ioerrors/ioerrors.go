package ioerrors

import (
	"io"

	"github.com/go-playground/errors"
)

// IOErrors helps classify io related errors
func IOErrors(c errors.Chain, err error) (cont bool) {
	switch err {
	case io.EOF:
		c.AddTypes("io")
		return
	case io.ErrClosedPipe:
		c.AddTypes("Permanent", "io")
		return
	case io.ErrNoProgress:
		c.AddTypes("Permanent", "io")
		return
	case io.ErrShortBuffer:
		c.AddTypes("Permanent", "io")
		return
	case io.ErrShortWrite:
		c.AddTypes("Permanent", "io")
		return
	case io.ErrUnexpectedEOF:
		c.AddTypes("Transient", "io")
		return
	}
	return true
}
