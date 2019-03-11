package errors

import (
	"runtime"
	"strings"
)

// Frame wraps a runtime.Frame to provide some helper functions while still allowing access to
// the original runtime.Frame
type Frame struct {
	runtime.Frame
}

// File is the runtime.Frame.File stripped down to just the filename
func (f Frame) File() string {
	name := f.Frame.File
	i := strings.LastIndexByte(name, '/')
	return name[i+1:]
}

// Line is the line of the runtime.Frame and exposed for convenience.
func (f Frame) Line() int {
	return f.Frame.Line
}

// Function is the runtime.Frame.Function stripped down to just the function name
func (f Frame) Function() string {
	name := f.Frame.Function
	i := strings.LastIndexByte(name, '.')
	return name[i+1:]
}

// Stack returns a stack from for parsing into a trace line
func Stack() Frame {
	return StackLevel(1)
}

// StackLevel returns a stack from for parsing into a trace line
// this is primarily used by other libraries who use this package
// internally as the level needs to be adjusted.
func StackLevel(skip int) (f Frame) {
	var frame [3]uintptr
	runtime.Callers(skip+2, frame[:])
	frames := runtime.CallersFrames(frame[:])
	f.Frame, _ = frames.Next()
	return
}
