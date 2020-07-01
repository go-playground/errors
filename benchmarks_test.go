package errors

import (
	"testing"
)

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

func BenchmarkErrorPrintWithTagsAndTypes(b *testing.B) {
	err := New("base error").AddTag("key", "value").AddTypes("Permanent", "Other")
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkErrorParallelPrintWithTagsAndTypes(b *testing.B) {
	err := New("base error").AddTag("key", "value").AddTypes("Permanent", "Other")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = err.Error()
		}
	})
}
