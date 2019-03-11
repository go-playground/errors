package neterrors

import (
	"net"

	"github.com/go-playground/errors"
)

const (
	permanent = "Permanent"
	transient = "Transient"
)

func init() {
	errors.RegisterHelper(NETErrors)
}

// NETErrors helps classify io related errors
func NETErrors(c errors.Chain, err error) (cont bool) {
	switch e := err.(type) {
	case *net.AddrError:
		tp := permanent
		if e.Temporary() {
			tp = transient
		}
		_ = c.AddTypes(tp, "net").AddTags(
			errors.T("addr", e.Addr),
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
		return false

	case *net.DNSError:
		tp := permanent
		if e.Temporary() {
			tp = transient
		}
		_ = c.AddTypes(tp, "net").AddTags(
			errors.T("name", e.Name),
			errors.T("server", e.Server),
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
		return false

	case *net.ParseError:
		_ = c.AddTypes(permanent, "net").AddTags(
			errors.T("type", e.Type),
			errors.T("text", e.Text),
		)
		return false

	case *net.OpError:
		tp := permanent
		if e.Temporary() {
			tp = transient
		}
		_ = c.AddTypes(tp, "net").AddTags(
			errors.T("op", e.Op),
			errors.T("net", e.Net),
			errors.T("addr", e.Addr),
			errors.T("local_addr", e.Source),
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
		return false
	case net.UnknownNetworkError:
		tp := permanent
		if e.Temporary() {
			tp = transient
		}
		_ = c.AddTypes(tp, "net").AddTags(
			errors.T("is_timeout", e.Timeout()),
			errors.T("is_temporary", e.Temporary()),
		)
	}

	switch err {
	case net.ErrWriteToConnected:
		_ = c.AddTypes(transient, "net")
		return false
	}
	return true
}
