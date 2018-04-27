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

func newLink(err error, prefix string, skipFrames int) *Link {
	return &Link{
		Err:    err,
		Prefix: prefix,
		Source: st(skipFrames),
	}

}

// Chain contains the chained errors, the links, of the chains if you will
type Chain []*Link

// Error returns the formatted error string
func (c Chain) Error() string {
	lines := make([]string, 0, len(c))
	//source=<source> <prefix>: <error> tag=value tag2=value2 types=type1,type2
	for i := len(c) - 1; i >= 0; i-- {
		line := c[i].formatError()
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
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
	Source string
}

// formatError prints a single Links error
func (l *Link) formatError() string {
	line := fmt.Sprintf("source=%s ", l.Source)

	if l.Prefix != "" {
		line += l.Prefix
	}

	if _, ok := l.Err.(Chain); !ok {
		if l.Prefix != "" {
			line += ": "
		}
		line += l.Err.Error()
	}

	for _, tag := range l.Tags {
		line += fmt.Sprintf(" %s=%v", tag.Key, tag.Value)
	}

	if len(l.Types) > 0 {
		line += " types=" + strings.Join(l.Types, ",")
	}
	return line
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
	return wrap(c, prefix, 0)
}
