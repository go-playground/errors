package errors

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	helpers []Helper
)

// RegisterHelper adds a new helper function to extract Type and Tag information.
// errors will run all registered helpers until a match is found.
// NOTE helpers are run in the order they are added.
func RegisterHelper(helper Helper) {
	for i := 0; i < len(helpers); i++ {
		if reflect.ValueOf(helpers[i]).Pointer() == reflect.ValueOf(helper).Pointer() {
			return
		}
	}
	helpers = append(helpers, helper)
}

// New creates an error with the provided text and automatically wraps it with line information.
func New(s string) Chain {
	return wrap(errors.New(s), "", 3)
}

// Newf creates an error with the provided text and automatically wraps it with line information.
// it also accepts a varadic for optional message formatting.
func Newf(format string, a ...interface{}) Chain {
	return wrap(fmt.Errorf(format, a...), "", 3)
}

// Wrap encapsulates the error, stores a contextual prefix and automatically obtains
// a stack trace.
func Wrap(err error, prefix string) Chain {
	return wrap(err, prefix, 3)
}

// Wrapf encapsulates the error, stores a contextual prefix and automatically obtains
// a stack trace.
// it also accepts a varadic for prefix formatting.
func Wrapf(err error, prefix string, a ...interface{}) Chain {
	return wrap(err, fmt.Sprintf(prefix, a...), 3)
}

// WrapSkipFrames is a special version of Wrap that skips extra n frames when determining error location.
// Normally only used when wrapping the library
func WrapSkipFrames(err error, prefix string, n uint) Chain {
	return wrap(err, prefix, int(n)+3)
}

func wrap(err error, prefix string, skipFrames int) (c Chain) {
	var ok bool
	if c, ok = err.(Chain); ok {
		c = append(c, newLink(err, prefix, skipFrames))
	} else {
		c = Chain{newLink(err, prefix, skipFrames)}
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
	}
	return false
}

// LookupTag recursively searches for the provided tag and returns it's value or nil
func LookupTag(err error, key string) interface{} {
	switch t := err.(type) {
	case Chain:
		for i := len(t) - 1; i >= 0; i-- {
			for j := 0; j < len(t[i].Tags); j++ {
				if t[i].Tags[j].Key == key {
					return t[i].Tags[j].Value
				}
			}
		}
	}
	return nil
}
