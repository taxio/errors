package errors

import (
	"errors"
	"runtime"
	"strings"
)

type stack []uintptr

type Error struct {
	message string
	err     error
	stack   stack
	attrs   map[string]any
}

func (e *Error) Error() string {
	var msgs []string

	if e.message != "" {
		msgs = append(msgs, e.message)
	}
	if e.err != nil {
		if len(e.err.Error()) > 0 {
			msgs = append(msgs, e.err.Error())
		}
	}
	return strings.Join(msgs, ": ")
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) StackTrace() []uintptr {
	return e.stack
}

func (e *Error) Attributes() map[string]any {
	return e.attrs
}

func New(message string, annotators ...any) error {
	err := &Error{
		message: message,
		err:     nil,
		stack:   callers(),
		attrs:   nil,
	}

	if len(annotators) == 0 {
		return err
	}

	return wrap(err, annotators...)
}

func Const(message string) error {
	return New(message, WithNoStackTrace())
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Join(errs ...error) error {
	return Wrap(errors.Join(errs...))
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Wrap(err error, annotators ...any) error {
	if err == nil {
		return nil
	}

	e := &Error{
		message: "",
		err:     err,
		stack:   callers(),
		attrs:   nil,
	}

	if x, ok := err.(interface{ Attributes() map[string]any }); ok {
		e.attrs = x.Attributes()
	}

	return wrap(e, annotators...)
}

func wrap(err *Error, annotators ...any) error {
	for _, annotator := range annotators {
		switch a := annotator.(type) {
		case string:
			WithMessage(a)(err)
		case Attribute:
			WithAttrs(a)(err)
		case AnnotatorFunc:
			a(err)
		default:
			// Do nothing (or should I panic?)
		}
	}
	return err
}

func callers() stack {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}
