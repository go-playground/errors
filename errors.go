package errors

import (
	"errors"
	stderrors "errors"
	"fmt"
	"reflect"
)

type unwrap interface{ Unwrap() error }
type is interface{ Is(error) bool }

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
	for {
		switch t := err.(type) {
		case Chain: // fast path
			err = t[0].Err
			continue
		case unwrap:
			if unwrappedErr := t.Unwrap(); unwrappedErr != nil {
				err = unwrappedErr
				continue
			}
		}
		return err
	}
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

// Is is to allow this library to be a drop-in replacement to the std library.
// I highly recommend using a switch T:= err.(type) instead.
//
// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == os.ErrExist }
//
// then Is(MyError{}, os.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
	//if target == nil {
	//	return err == target
	//}
	//
	//isComparable := reflect.TypeOf(target).Comparable()
	//for {
	//	if isComparable && err == target {
	//		return true
	//	}
	//	if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
	//		return true
	//	}
	//	if err = stderrors.Unwrap(err); err == nil {
	//		return false
	//	}
	//}
}
