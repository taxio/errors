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

func New(message string) error {
	return &Error{
		message: message,
		err:     nil,
		stack:   callers(),
		attrs:   nil,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Wrap(err error, annotators ...Annotator) error {
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

	for _, annotator := range annotators {
		annotator(e)
	}

	return e
}

func callers() stack {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}
