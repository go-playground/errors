package errors

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"

	runtimeext "github.com/go-playground/pkg/v5/runtime"
	unsafeext "github.com/go-playground/pkg/v5/unsafe"
)

var (
	_ unwrap = (*Chain)(nil)
	_ is     = (*Chain)(nil)
	_ as     = (*Chain)(nil)
)

// T is a shortcut to make a Tag
func T(key string, value interface{}) Tag {
	return Tag{Key: key, Value: value}
}

// Tag contains a single key value combination
// to be attached to your error
type Tag struct {
	Key   string
	Value interface{}
}

func newLink(err error, prefix string, skipFrames int) *Link {
	return &Link{
		Err:    err,
		Prefix: prefix,
		Source: runtimeext.StackLevel(skipFrames),
	}

}

// Chain contains the chained errors, the links, of the chains if you will
type Chain []*Link

// Error returns the formatted error string
func (c Chain) Error() string {
	//if errFormatFn != nil {
	return errFormatFn(c)
	//}
	//b := make([]byte, 0, len(c)*192)
	//
	//for i := 0; i < len(c); i++ {
	//	b = c[i].formatError(b)
	//	b = append(b, '\n')
	//}
	//return unsafeext.BytesToString(b[:len(b)-1])
}

// Link contains a single error entry, unless it's the top level error, in
// which case it only contains an array of errors
type Link struct {

	// Err is the wrapped error, either the original or already wrapped
	Err error

	// Prefix contains the error prefix text
	Prefix string

	// Type stores one or more categorized types of error set by the caller using AddTypes and is optional
	Types []string

	// Tags contains an array of tags associated with this error, if any
	Tags []Tag

	// Source contains the name, file and lines obtained from the stack trace
	Source runtimeext.Frame
}

// formatError prints a single Links error
func (l *Link) formatError(b []byte) []byte {
	b = append(b, "source="...)
	idx := strings.LastIndexByte(l.Source.Frame.Function, '.')
	if idx == -1 {
		b = append(b, l.Source.File()...)
	} else {
		b = append(b, l.Source.Frame.Function[:idx]...)
		b = append(b, '/')
		b = append(b, l.Source.File()...)
	}
	b = append(b, ':')
	b = strconv.AppendInt(b, int64(l.Source.Line()), 10)
	b = append(b, ':')
	b = append(b, l.Source.Frame.Function[idx+1:]...)
	b = append(b, ' ')
	b = append(b, "error="...)

	if l.Prefix != "" {
		b = append(b, l.Prefix...)
	}

	if _, ok := l.Err.(Chain); !ok {
		if l.Prefix != "" {
			b = append(b, ": "...)
		}
		b = append(b, l.Err.Error()...)
	}

	for _, tag := range l.Tags {
		b = append(b, ' ')
		b = append(b, tag.Key...)
		b = append(b, '=')
		switch t := tag.Value.(type) {
		case string:
			b = append(b, t...)
		case int:
			b = strconv.AppendInt(b, int64(t), 10)
		case int8:
			b = strconv.AppendInt(b, int64(t), 10)
		case int16:
			b = strconv.AppendInt(b, int64(t), 10)
		case int32:
			b = strconv.AppendInt(b, int64(t), 10)
		case int64:
			b = strconv.AppendInt(b, t, 10)
		case uint:
			b = strconv.AppendUint(b, uint64(t), 10)
		case uint8:
			b = strconv.AppendUint(b, uint64(t), 10)
		case uint16:
			b = strconv.AppendUint(b, uint64(t), 10)
		case uint32:
			b = strconv.AppendUint(b, uint64(t), 10)
		case uint64:
			b = strconv.AppendUint(b, t, 10)
		case float32:
			b = strconv.AppendFloat(b, float64(t), 'g', -1, 32)
		case float64:
			b = strconv.AppendFloat(b, t, 'g', -1, 64)
		case bool:
			b = strconv.AppendBool(b, t)
		default:
			b = append(b, fmt.Sprintf("%v", tag.Value)...)
		}
	}

	if len(l.Types) > 0 {
		b = append(b, " types="...)
		for _, t := range l.Types {
			b = append(b, t...)
			b = append(b, ',')
		}
		b = b[:len(b)-1]
	}
	return b
}

// helper method to get the current *Link from the top level
func (c Chain) current() *Link {
	return c[len(c)-1]
}

// AddTags allows the addition of multiple tags
func (c Chain) AddTags(tags ...Tag) Chain {
	l := c.current()
	if len(l.Tags) == 0 {
		l.Tags = make([]Tag, 0, len(tags))
	}
	l.Tags = append(l.Tags, tags...)
	return c
}

// AddTag allows the addition of a single tag
func (c Chain) AddTag(key string, value interface{}) Chain {
	return c.AddTags(Tag{Key: key, Value: value})
}

// AddTypes sets one or more categorized types on the Link error
func (c Chain) AddTypes(typ ...string) Chain {
	l := c.current()
	l.Types = append(l.Types, typ...)
	return c
}

// Wrap adds another contextual prefix to the error chain
func (c Chain) Wrap(prefix string) Chain {
	return wrap(c, prefix, 3)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// If attempting to retrieve the cause see Cause function instead.
func (c Chain) Unwrap() error {
	if len(c) == 1 {
		if e, ok := c[0].Err.(unwrap); ok {
			return e.Unwrap()
		}
		return c[0].Err
	}
	return c[:len(c)-1]
}

// Is reports whether any error in err's chain matches target.
//
// This is here to help make it a drop in replacement to the std error handler, I highly recommend using
// Cause + switch statement instead.
func (c Chain) Is(target error) bool {
	return stderrors.Is(c[len(c)-1].Err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// a different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
func (c Chain) As(target interface{}) bool {
	return stderrors.As(c[len(c)-1].Err, target)
}

func defaultFormatFn(c Chain) string {
	b := make([]byte, 0, len(c)*192)

	for i := 0; i < len(c); i++ {
		b = c[i].formatError(b)
		b = append(b, '\n')
	}
	return unsafeext.BytesToString(b[:len(b)-1])
}
