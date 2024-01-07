package errors

import "runtime"

type Annotator func(*Error)

func callers() *stack {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func WithMessage(message string) Annotator {
	return func(err *Error) {
		err.message = message
	}
}

func WithStack() Annotator {
	return func(err *Error) {
		err.stack = callers()
	}
}

func WithAttrs(attrs map[string]any) Annotator {
	return func(err *Error) {
		err.attrs = attrs
	}
}
