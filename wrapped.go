package errors

import (
	"fmt"
	"strings"
)

// T is a shortcut to make a Tag
func T(key string, value interface{}) Tag {
	return Tag{Key: key, Value: value}
}

// Tag contains a single key value conbination
// to be attached to your error
type Tag struct {
	Key   string
	Value interface{}
}

func newWrapped(err error, prefix string) *Wrapped {
	return &Wrapped{
		Err:    err,
		Prefix: prefix,
		Source: st(),
	}
}

// Wrapped contains a single error entry, unless it's the top level error, in
// which case it only contains an array of errors
type Wrapped struct {
	// hidden field with wrapped errors, will expose a helper method to get this
	Errors []*Wrapped

	// Err is the wrapped error, either the original or already wrapped
	Err error

	// Prefix contains the error prefix text
	Prefix string

	// Type stores one or more categorized types of error set by the caller using WithTypes and is optional
	Types []string

	// Tags contains an array of tags associated with this error, if any
	Tags []Tag

	// Source contains the name, file and lines obtained from the stack trace
	Source string
}

// Error returns the formatted error string
func (w *Wrapped) Error() string {
	if len(w.Errors) > 0 {
		lines := make([]string, 0, len(w.Errors))
		// source=<source> <prefix>: <error> tag=value tag2=value2 types=type1,type2
		for i := len(w.Errors) - 1; i >= 0; i-- {
			line := w.Errors[i].parseLine()
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n")
	}
	return w.parseLine()
}

func (w *Wrapped) parseLine() string {
	line := fmt.Sprintf("source=%s %s", w.Source, w.Prefix)

	if _, isWrapped := w.Err.(*Wrapped); !isWrapped {
		line += ": " + w.Err.Error()
	}

	for _, tag := range w.Tags {
		line += fmt.Sprintf(" %s=%v", tag.Key, tag.Value)
	}

	if len(w.Types) > 0 {
		line += " types=" + strings.Join(w.Types, ",")
	}
	return line
}

// helper method to get the current *Wrapped from the top level
func (w *Wrapped) current() *Wrapped {
	return w.Errors[len(w.Errors)-1]
}

// WithTags allows the addition of multiple tags
func (w *Wrapped) WithTags(tags ...Tag) *Wrapped {
	wr := w.current()
	if len(wr.Tags) == 0 {
		wr.Tags = make([]Tag, 0, len(tags))
	}
	wr.Tags = append(wr.Tags, tags...)
	return w
}

// WithTag allows the addition of a single tag
func (w *Wrapped) WithTag(key string, value interface{}) *Wrapped {
	return w.WithTags(Tag{Key: key, Value: value})
}

// WithTypes sets one or more categorized types on the Wrapped error
func (w *Wrapped) WithTypes(typ ...string) *Wrapped {
	wr := w.current()
	wr.Types = append(wr.Types, typ...)
	return w
}
