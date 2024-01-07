package errors_test

import (
	stderr "errors"
	"testing"

	"github.com/taxio/errors"
)

func TestNew(t *testing.T) {
	msg := "test"
	err := errors.New(msg)
	if err.Error() != msg {
		t.Errorf("expected %s, got %s", msg, err.Error())
	}
}

func TestIs(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		err := errors.New("test")
		if !errors.Is(err, err) {
			t.Errorf("expected true, got false")
		}
	})
	t.Run("different", func(t *testing.T) {
		err1 := errors.New("test")
		err2 := errors.New("test")
		if errors.Is(err1, err2) {
			t.Errorf("expected false, got true")
		}
	})
}

func TestAs(t *testing.T) {
	t.Skip("TODO")
}

func TestWrap(t *testing.T) {
	t.Skip("TODO")
}

func TestWithMessage(t *testing.T) {
	t.Skip("TODO")
}

func TestWithStack(t *testing.T) {
	t.Skip("TODO")
}

func TestWithAttrs(t *testing.T) {
	t.Skip("TODO")
}

func TestUnwrap(t *testing.T) {
	t.Skip("TODO")
}

func TestJoin(t *testing.T) {
	tests := map[string]struct {
		errs []error
		want string
	}{
		"empty": {
			errs: []error{},
			want: "",
		},
		"single": {
			errs: []error{errors.New("test")},
			want: "test",
		},
		"single w/ std errors": {
			errs: []error{stderr.New("test")},
			want: "test",
		},
		"2 errors": {
			errs: []error{errors.New("test1"), errors.New("test2")},
			want: "test1\ntest2",
		},
		"3 errors": {
			errs: []error{errors.New("test1"), errors.New("test2"), errors.New("test3")},
			want: "test1\ntest2\ntest3",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			err := errors.Join(tt.errs...)
			if err == nil {
				if tt.want != "" {
					t.Errorf("expected %s, got nil", tt.want)
				}
				return
			}
			if err.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, err.Error())
			}
		})
	}
}
