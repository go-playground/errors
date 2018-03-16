package errors

// Wrap encapsulates the error, stores a contextual prefix and automatically obtains
// a stack trace.
func Wrap(err error, prefix string) *Wrapped {
	if w, ok := err.(*Wrapped); ok {
		w.Errors = append(w.Errors, newWrapped(err, prefix))
		return w
	}
	return &Wrapped{
		Errors: []*Wrapped{newWrapped(err, prefix)},
	}
}

// HasType is a helper function that will recurse up from the root error and check that the provided type
// is present using an equality check
func HasType(err error, typ string) bool {
	w, ok := err.(*Wrapped)
	if !ok {
		return false
	}
	for i := len(w.Errors) - 1; i >= 0; i-- {
		for j := 0; j < len(w.Errors[i].Types); j++ {
			if w.Errors[i].Types[j] == typ {
				return true
			}
		}
	}
	return false
}

// Cause extracts and returns the root error
func Cause(err error) error {
	if w, ok := err.(*Wrapped); ok {
		// if root level error
		if len(w.Errors) > 0 {
			return w.Errors[0]
		}
		// already extracted error
		return w
	}
	return err
}

// IsErr will fetch the root error, and check the original error against the provided type
// eg. errors.IsErr(io.EOF)
func IsErr(err, errType error) bool {
	if w, ok := err.(*Wrapped); ok {
		// if root level error
		if len(w.Errors) > 0 {
			return w.Errors[0].Err == errType
		}
		// already extracted error
		return w.Err == errType
	}
	return err == errType
}
