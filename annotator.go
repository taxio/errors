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

type Attribute struct {
	key   string
	value any
}

func Attr(key string, value any) Attribute {
	return Attribute{key: key, value: value}
}

func WithAttrs(attrs ...Attribute) Annotator {
	return func(err *Error) {
		if err.attrs == nil {
			err.attrs = make(map[string]any, len(attrs))
		}
		for _, attr := range attrs {
			err.attrs[attr.key] = attr.value
		}
	}
}
