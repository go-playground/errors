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
	assertError(t, err)
	err = Wrap(defaultErr, "prefix")
	assertError(t, err)
	err = Wrap(err, "prefix2")
	assertError(t, err)
	err = fmt.Errorf("wrapping Chain: %w", err)
	assertError(t, err)
	err = fmt.Errorf("wrapping err again: %w", err)
	assertError(t, err)
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

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("err is nil")
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
		name string
		typ  string
		err  error
	}{
		{
			name: "basic types",
			typ:  "Permanent",
			err:  Wrap(fmt.Errorf("this is an %s", "error"), "prefix").AddTypes("Permanent", "internalError"),
		},
		{
			name: "std wrapped",
			typ:  "MyType",
			err:  fmt.Errorf("std wrapped %w", New("base error").AddTypes("MyType")),
		},
	}

	for i, tc := range tests {
		tc := tc
		i := i
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !HasType(tc.err, tc.typ) {
				t.Fatalf("IDX: %d want %t got %t", i, true, false)
			}
		})
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
	assertError(t, err)
	err = Wrap(defaultErr, "prefix")
	assertError(t, err)
	err = Wrap(err, "prefix2")
	assertError(t, err)
	err = fmt.Errorf("wrapping Chain: %w", err)
	assertError(t, err)
	err = fmt.Errorf("wrapping err again: %w", err)
	assertError(t, err)
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
	key := "Key"
	value := "Value"

	tests := []struct {
		name  string
		err   error
		key   string
		value any
	}{
		{
			name: "basic wrap",
			err:  Wrap(io.EOF, "prefix").AddTag(key, value),
		},
		{
			name: "double wrapped",
			err:  Wrap(Wrap(io.EOF, "prefix").AddTag(key, value), "wrapped"),
		},
		{
			name: "std lib wrapped",
			err:  fmt.Errorf("wrapped %w", Wrap(io.EOF, "prefix").AddTag(key, value)),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if LookupTag(tc.err, key) != value {
				t.Fatalf("want '%s' got '%v'", value, LookupTag(tc.err, key))
			}
		})
	}
}

func TestIs(t *testing.T) {
	defaultErr := io.EOF
	err := fmt.Errorf("std wrapped: %w", defaultErr)
	assertError(t, err)
	err = Wrap(defaultErr, "prefix")
	assertError(t, err)
	err = Wrap(err, "prefix2")
	assertError(t, err)
	err = fmt.Errorf("wrapping Chain: %w", err)
	assertError(t, err)
	err = fmt.Errorf("wrapping err again: %w", err)
	assertError(t, err)
	err = Wrap(err, "wrapping std with chain")
	if !Is(err, io.EOF) {
		t.Fatal("want true got false")
	}
}

type myErrorType struct {
	msg string
}

func (e *myErrorType) Error() string {
	return e.msg
}

func TestAs(t *testing.T) {
	defaultErr := &myErrorType{msg: "my error type"}
	err := fmt.Errorf("std wrapped: %w", defaultErr)
	assertError(t, err)
	err = Wrap(defaultErr, "prefix")
	assertError(t, err)
	err = Wrap(err, "prefix2")
	assertError(t, err)
	err = fmt.Errorf("wrapping Chain: %w", err)
	assertError(t, err)
	err = fmt.Errorf("wrapping err again: %w", err)
	assertError(t, err)
	err = Wrap(err, "wrapping std with chain")

	var myErr *myErrorType
	if !As(err, &myErr) {
		t.Fatal("want true got false")
	}
}

func TestCustomFormatFn(t *testing.T) {
	RegisterErrorFormatFn(func(c Chain) (s string) {
		return c[0].Err.Error()
	})
	err := io.EOF
	err = Wrap(err, "my error prefix")
	if err.Error() != "EOF" {
		t.Errorf("Expected output of 'EOF'")
	}
}
