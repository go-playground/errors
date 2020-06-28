package errors

import (
	stderrors "errors"
	"fmt"
	"testing"
)

func BenchmarkErrorStd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Errorf("blah %w", stderrors.New("base error"))
	}
}

func BenchmarkErrorNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New("base error")
	}
}

func BenchmarkErrorParallelNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New("base error")
		}
	})
}

func BenchmarkErrorPrint(b *testing.B) {
	err := New("base error")
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkErrorParallelPrint(b *testing.B) {
	err := New("base error")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = err.Error()
		}
	})
}
