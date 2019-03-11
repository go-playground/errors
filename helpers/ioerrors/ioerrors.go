package ioerrors

import (
	"io"

	"github.com/go-playground/errors"
)

func init() {
	errors.RegisterHelper(IOErrors)
}

// IOErrors helps classify io related errors
func IOErrors(c errors.Chain, err error) (cont bool) {
	switch err {
	case io.EOF:
		_ = c.AddTypes("io")
		return
	case io.ErrClosedPipe:
		_ = c.AddTypes("Permanent", "io")
		return
	case io.ErrNoProgress:
		_ = c.AddTypes("Permanent", "io")
		return
	case io.ErrShortBuffer:
		_ = c.AddTypes("Permanent", "io")
		return
	case io.ErrShortWrite:
		_ = c.AddTypes("Permanent", "io")
		return
	case io.ErrUnexpectedEOF:
		_ = c.AddTypes("Transient", "io")
		return
	}
	return true
}
