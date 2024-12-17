package errors

type AnnotatorFunc func(*Error)

func WithMessage(message string) AnnotatorFunc {
	return func(err *Error) {
		if message == "" {
			return
		}

		if err.message == "" {
			err.message = message
		} else {
			err.message = message + ": " + err.message
		}
	}
}

type Attribute struct {
	key   string
	value any
}

func Attr(key string, value any) Attribute {
	return Attribute{key: key, value: value}
}

func WithAttrs(attrs ...Attribute) AnnotatorFunc {
	return func(err *Error) {
		if err.attrs == nil {
			err.attrs = make(map[string]any, len(attrs))
		}
		for _, attr := range attrs {
			err.attrs[attr.key] = attr.value
		}
	}
}

func WithNoStackTrace() AnnotatorFunc {
	return func(err *Error) {
		err.stack = nil
	}
}
