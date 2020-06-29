package errors

import (
	stderrors "errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	defaultErr := fmt.Errorf("this is an %s", "error")
	testWrapper := func(err error, prefix string) Chain {
		return WrapSkipFrames(err, prefix, 1)
	}

	err0 := New("42")
	err1 := Wrap(defaultErr, "prefix 1")
	err2 := err1.Wrap("prefix 2")
	err3 := testWrapper(err2, "prefix 3")
	err4 := Newf("this is an %s", "error")
	err5 := Wrapf(defaultErr, "this is an %s", "error")

	tests := []struct {
		err Chain
		pre string
		suf string
	}{
		{
			err: err0,
			pre: "TestWrap: ",
			suf: "errors_test.go:17",
		},
		{
			err: err1,
			pre: "TestWrap: ",
			suf: "errors_test.go:18",
		},
		{
			err: err2,
			pre: "TestWrap: ",
			suf: "errors_test.go:19",
		},
		{
			err: err3,
			pre: "TestWrap: ",
			suf: "errors_test.go:20",
		},
		{
			err: err4,
			pre: "TestWrap: ",
			suf: "errors_test.go:21",
		},
		{
			err: err5,
			pre: "TestWrap: ",
			suf: "errors_test.go:22",
		},
	}

	for i, tt := range tests {
		link := tt.err.current()
		source := fmt.Sprintf("%s: %s:%d", link.Source.Function(), link.Source.File(), link.Source.Line())
		if !strings.HasSuffix(source, tt.suf) || !strings.HasPrefix(source, tt.pre) {
			t.Fatalf("IDX: %d want %s<path>%s got %s", i, tt.pre, tt.suf, source)
		}
	}
}

func TestUnwrap(t *testing.T) {
	defaultErr := stderrors.New("this is an error")
	err := fmt.Errorf("std wrapped: %w", defaultErr)
	err = Wrap(defaultErr, "prefix")
	err = Wrap(err, "prefix2")
	err = fmt.Errorf("wrapping Chain: %w", err)
	err = fmt.Errorf("wrapping err again: %w", err)
	err = Wrap(err, "wrapping std with chain")

	for {
		switch t := err.(type) {
		case unwrap:
			if unwrappedErr := t.Unwrap(); unwrappedErr != nil {
				err = unwrappedErr
				continue
			}
		}
		break
	}
	expect := defaultErr.Error()
	if err.Error() != expect {
		t.Fatalf("want %s got %s", expect, err.Error())
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
		t.Fatalf("want %s got %s", expect, cause.Error())
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

func TestCauseStdErrorsMixed(t *testing.T) {
	defaultErr := stderrors.New("this is an error")
	err := fmt.Errorf("std wrapped: %w", defaultErr)
	err = Wrap(defaultErr, "prefix")
	err = Wrap(err, "prefix2")
	err = fmt.Errorf("wrapping Chain: %w", err)
	err = fmt.Errorf("wrapping err again: %w", err)
	err = Wrap(err, "wrapping std with chain")
	cause := Cause(err)
	expect := defaultErr.Error()
	if cause.Error() != expect {
		t.Fatalf("want %s got %s", expect, cause.Error())
	}
}

func TestHelpers(t *testing.T) {
	fn := func(w Chain, _ error) (cont bool) {
		_ = w.AddTypes("Test").AddTags(T("test", "tag")).AddTag("foo", "bar")
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

func TestIs(t *testing.T) {
	defaultErr := io.EOF
	err := fmt.Errorf("std wrapped: %w", defaultErr)
	err = Wrap(defaultErr, "prefix")
	err = Wrap(err, "prefix2")
	err = fmt.Errorf("wrapping Chain: %w", err)
	err = fmt.Errorf("wrapping err again: %w", err)
	err = Wrap(err, "wrapping std with chain")
	if !Is(err, io.EOF) {
		t.Fatal("want true got false")
	}
}
