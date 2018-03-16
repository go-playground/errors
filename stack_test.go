package errors

import (
	"fmt"
	"strings"
	"testing"
)

func TestStack(t *testing.T) {
	tests := []struct {
		file string
		line string
		name string
	}{
		{
			file: "stack_test.go",
			line: "23",
			name: "TestStack",
		},
	}

	for i, tt := range tests {
		frame := Stack()
		if v := fmt.Sprintf("%n", frame); tt.name != v {
			t.Fatalf("IDX: %d want %s got %s", i, tt.name, v)
		}
		if v := fmt.Sprintf("%+s", frame); !strings.HasSuffix(v, tt.file) {
			t.Fatalf("IDX: %d want %s got %s", i, tt.file, v)
		}
		if v := fmt.Sprintf("%d", frame); tt.line != v {
			t.Fatalf("IDX: %d want %s got %s", i, tt.line, v)
		}
	}
}
