package neterrors

import (
	"net"

	"github.com/go-playground/errors"
)

// NETErrors helps classify io related errors
func NETErrors(w *errors.Wrapped, err error) (cont bool) {
	switch e := err.(type) {
	case *net.AddrError:
		tp := "Permanent"
		if e.Temporary() {
			tp = "Transient"
		}
		w.WithTypes(tp, "net").WithTags(
			errors.T("addr", e.Addr),
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
		return false

	case *net.DNSError:
		tp := "Permanent"
		if e.Temporary() {
			tp = "Transient"
		}
		w.WithTypes(tp, "net").WithTags(
			errors.T("name", e.Name),
			errors.T("server", e.Server),
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
		return false

	case *net.ParseError:
		w.WithTypes("Permanent", "net").WithTags(
			errors.T("type", e.Type),
			errors.T("text", e.Text),
		)
		return false

	case *net.OpError:
		tp := "Permanent"
		if e.Temporary() {
			tp = "Transient"
		}
		w.WithTypes(tp, "net").WithTags(
			errors.T("op", e.Op),
			errors.T("net", e.Net),
			errors.T("addr", e.Addr),
			errors.T("local_addr", e.Source),
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
		return false
	case net.UnknownNetworkError:
		tp := "Permanent"
		if e.Temporary() {
			tp = "Transient"
		}
		w.WithTypes(tp, "net").WithTags(
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
	}

	switch err {
	case net.ErrWriteToConnected:
		w.WithTypes("Transient", "net")
		return false
	}
	return true
}
