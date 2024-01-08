package errors

type Annotator func(*Error)

func WithMessage(message string) Annotator {
	return func(err *Error) {
		err.message = message
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
