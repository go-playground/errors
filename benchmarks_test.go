package errors

import (
	"testing"
)

func BenchmarkError(b *testing.B) {
	err := New("base error")
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkErrorParallel(b *testing.B) {
	err := New("base error")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = err.Error()
		}
	})
}
