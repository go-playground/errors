package errors

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	defaultErr := fmt.Errorf("this is an %s", "error")
	testWrapper := func (err error, prefix string) Chain {
		return WrapSkipFrames(err, prefix, 1)
	}

	err0 := New("42")
	err1 := Wrap(defaultErr, "prefix 1")
	err2 := err1.Wrap("prefix 2")
	err3 := testWrapper(err2, "prefix 3")

	tests := []struct {
		err Chain
		pre string
		suf string
	}{
		{
			err: err0,
			pre: "TestWrap: ",
			suf: "errors_test.go:16",
		},
		{
			err: err1,
			pre: "TestWrap: ",
			suf: "errors_test.go:17",
		},
		{
			err: err2,
			pre: "TestWrap: ",
			suf: "errors_test.go:18",
		},
		{
			err: err3,
			pre: "TestWrap: ",
			suf: "errors_test.go:19",
		},
	}

	for i, tt := range tests {
		link := tt.err.current()
		if !strings.HasSuffix(link.Source, tt.suf) || !strings.HasPrefix(link.Source, tt.pre) {
			t.Fatalf("IDX: %d want %s<path>%s got %s", i, tt.pre, tt.suf, link.Source)
		}
	}
}

func TestTags(t *testing.T) {
	tests := []struct {
		err  string
		tags []Tag
	}{
		{
			err:  "key=value key2=value2",
			tags: []Tag{T("key", "value"), T("key2", "value2")},
		},
	}

	defaultErr := fmt.Errorf("this is an %s", "error")

	for i, tt := range tests {
		err := Wrap(defaultErr, "prefix").AddTags(tt.tags...)
		if !strings.HasSuffix(err.Error(), tt.err) {
			t.Fatalf("IDX: %d want %s got %s", i, tt.err, err.Error())
		}
	}
}

func TestTypes(t *testing.T) {
	tests := []struct {
		err   string
		tags  []Tag
		types []string
	}{
		{
			err:   "types=Permanent,InternalError",
			tags:  []Tag{T("key", "value"), T("key2", "value2")},
			types: []string{"Permanent", "InternalError"},
		},
	}

	defaultErr := fmt.Errorf("this is an %s", "error")

	for i, tt := range tests {
		err := Wrap(defaultErr, "prefix").AddTags(tt.tags...).AddTypes(tt.types...)
		if !strings.HasSuffix(err.Error(), tt.err) {
			t.Fatalf("IDX: %d want %s got %s", i, tt.err, err.Error())
		}
	}
}

func TestHasType(t *testing.T) {
	tests := []struct {
		types []string
		typ   string
	}{
		{
			types: []string{"Permanent", "internalError"},
			typ:   "Permanent",
		},
	}

	defaultErr := fmt.Errorf("this is an %s", "error")

	for i, tt := range tests {
		err := Wrap(defaultErr, "prefix").AddTypes(tt.types...)
		if !HasType(err, tt.typ) {
			t.Fatalf("IDX: %d want %t got %t", i, true, false)
		}
	}
}

func TestCause(t *testing.T) {
	defaultErr := fmt.Errorf("this is an %s", "error")
	err := Wrap(defaultErr, "prefix")
	err = Wrap(err, "prefix2")
	cause := Cause(err)
	expect := "this is an error"
	if cause.Error() != expect {
		t.Fatalf("want %s got %s", expect, err.Error())
	}
}

func TestCause2(t *testing.T) {
	err := Wrap(io.EOF, "prefix")
	err = Wrap(err, "prefix2")

	if Cause(err) != io.EOF {
		t.Fatalf("want %t got %t", true, false)
	}
	cause := Cause(err)
	if Cause(cause) != io.EOF {
		t.Fatalf("want %t got %t", true, false)
	}
	if Cause(io.EOF) != io.EOF {
		t.Fatalf("want %t got %t", true, false)
	}
}

func TestHelpers(t *testing.T) {
	fn := func(w Chain, err error) (cont bool) {
		w.AddTypes("Test").AddTags(T("test", "tag")).AddTag("foo", "bar")
		return false
	}
	RegisterHelper(fn)

	err := Wrap(io.EOF, "prefix")
	if !HasType(err, "Test") {
		t.Errorf("Expected to have type 'Test'")
	}
}

func TestLookupTag(t *testing.T) {
	err := Wrap(io.EOF, "prefix").AddTag("Key", "Value")
	if LookupTag(err, "Key").(string) != "Value" {
		t.Fatalf("want %s got %v", "Value", LookupTag(err, "Key"))
	}
}
