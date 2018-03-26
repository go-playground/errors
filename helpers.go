package errors

// Helper is a function which will automatically extract Type and Tag information based on the supplied err and
// add it to the supplied *Wrapped error; this can be used independently or by registering using errors.REgisterHelper(...),
// which will run the registered helper every time errors.Wrap(...) is called.
type Helper func(*Wrapped, error) bool
