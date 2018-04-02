package errors

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		pre string
		suf string
	}{
		{
			pre: "source=TestWrap: ",
			suf: "errors_test.go:24 prefix: this is an error",
		},
	}

	defaultErr := fmt.Errorf("this is an %s", "error")

	for i, tt := range tests {
		err := Wrap(defaultErr, "prefix")
		if !strings.HasSuffix(err.Error(), tt.suf) || !strings.HasPrefix(err.Error(), tt.pre) {
			t.Fatalf("IDX: %d want %s<path>%s got %s", i, tt.pre, tt.suf, err.Error())
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
		err := Wrap(defaultErr, "prefix").WithTags(tt.tags...)
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
		err := Wrap(defaultErr, "prefix").WithTags(tt.tags...).WithTypes(tt.types...)
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
		err := Wrap(defaultErr, "prefix").WithTypes(tt.types...)
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

func TestIsErr(t *testing.T) {
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
		w.WithTypes("Test").WithTags(T("test", "tag")).WithTag("foo", "bar")
		return false
	}
	RegisterHelper(fn)

	err := Wrap(io.EOF, "prefix")
	if !HasType(err, "Test") {
		t.Errorf("Expected to have type 'Test'")
	}
}
