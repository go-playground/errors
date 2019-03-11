package errors

import (
	"testing"
)

func TestStack(t *testing.T) {
	tests := []struct {
		file string
		line int
		fn   string
	}{
		{
			file: "stack_test.go",
			line: 21,
			fn:   "TestStack",
		},
	}

	for i, tt := range tests {
		frame := Stack()
		if tt.fn != frame.Function() {
			t.Fatalf("IDX: %d want %s got %s", i, tt.fn, frame.Function())
		}
		if tt.file != frame.File() {
			t.Fatalf("IDX: %d want %s got %s", i, tt.file, frame.File())
		}
		if tt.line != frame.Line() {
			t.Fatalf("IDX: %d want %d got %d", i, tt.line, frame.Line())
		}
	}
}
