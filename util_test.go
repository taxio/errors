package errors_test

import (
	stderr "errors"
	"testing"

	"github.com/taxio/errors"
)

func TestBaseStackTrace(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got := errors.BaseStackTrace(nil)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("other error", func(t *testing.T) {
		got := errors.BaseStackTrace(stderr.New("test"))
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("single error", func(t *testing.T) {
		err := errors.New("test")
		cErr := mustCast(t, err)

		got := errors.BaseStackTrace(err)
		if len(cErr.StackTrace()) != len(got) {
			t.Fatalf("expected %d, got %d", len(cErr.StackTrace()), len(got))
		}
		for i, st := range cErr.StackTrace() {
			if got[i] != st {
				t.Errorf("expected %v, got %v", st, got[i])
			}
		}
	})

	t.Run("nested", func(t *testing.T) {
		baseErr := errors.New("base")
		cBaseErr := mustCast(t, baseErr)
		wrapErr1 := errors.Wrap(baseErr)
		wrapErr2 := errors.Wrap(wrapErr1)

		got := errors.BaseStackTrace(wrapErr2)
		if len(cBaseErr.StackTrace()) != len(got) {
			t.Fatalf("expected %d, got %d", len(cBaseErr.StackTrace()), len(got))
		}
		for i, st := range cBaseErr.StackTrace() {
			if got[i] != st {
				t.Errorf("expected %v, got %v", st, got[i])
			}
		}
	})

	t.Run("nested with other base error", func(t *testing.T) {
		baseErr := stderr.New("base")
		wrapErr1 := errors.Wrap(baseErr)
		cWrapErr1 := mustCast(t, wrapErr1)
		wrapErr2 := errors.Wrap(wrapErr1)

		got := errors.BaseStackTrace(wrapErr2)
		if len(cWrapErr1.StackTrace()) != len(got) {
			t.Fatalf("expected %d, got %d", len(cWrapErr1.StackTrace()), len(got))
		}
		for i, st := range cWrapErr1.StackTrace() {
			if got[i] != st {
				t.Errorf("expected %v, got %v", st, got[i])
			}
		}
	})
}

func mustCast(t *testing.T, err error) *errors.Error {
	t.Helper()
	var cErr *errors.Error
	if !errors.As(err, &cErr) {
		t.Fatalf("expected error to be *errors.Error")
	}
	return cErr
}
