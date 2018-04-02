package errors

import (
	"errors"
)

var (
	helpers []Helper
)

// RegisterHelper adds a new helper function to extract Type and Tag information.
// errors will run all registered helpers until a match is found.
// NOTE helpers are run in the order they are added.
func RegisterHelper(helper Helper) {
	helpers = append(helpers, helper)
}

// New creates an error with the provided text and automatically wraps it with line information.
func New(s string) Chain {
	return wrap(errors.New(s), "")
}

// Wrap encapsulates the error, stores a contextual prefix and automatically obtains
// a stack trace.
func Wrap(err error, prefix string) Chain {
	return wrap(err, prefix)
}

func wrap(err error, prefix string) (c Chain) {
	var ok bool
	if c, ok = err.(Chain); ok {
		c = append(c, newLink(err, prefix))
	} else {
		c = Chain{newLink(err, prefix)}
	}
	for _, h := range helpers {
		if !h(c, err) {
			break
		}
	}
	return
}

// Cause extracts and returns the root wrapped error (the naked error with no additional information
func Cause(err error) error {
	switch t := err.(type) {
	case Chain:
		return t[0].Err
	case *Link:
		return t.Err
	default:
		return err
	}
	// TODO: lookup via Cause interface recursively on error
}

// HasType is a helper function that will recurse up from the root error and check that the provided type
// is present using an equality check
func HasType(err error, typ string) bool {
	switch t := err.(type) {
	case Chain:
		for i := len(t) - 1; i >= 0; i-- {
			for j := 0; j < len(t[i].Types); j++ {
				if t[i].Types[j] == typ {
					return true
				}
			}
		}
	case *Link:
		for i := 0; i < len(t.Types); i++ {
			if t.Types[i] == typ {
				return true
			}
		}
	}
	return false
}
