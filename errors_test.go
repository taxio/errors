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

type CustomError struct {
	message string
}

func (e *CustomError) Error() string {
	return e.message
}

func TestAs(t *testing.T) {
	t.Run("same type", func(t *testing.T) {
		err := errors.New("test")
		var target *errors.Error
		if !errors.As(err, &target) {
			t.Errorf("expected true, got false")
		}
		if target == nil {
			t.Errorf("target is not set")
		}
	})

	t.Run("different type", func(t *testing.T) {
		err := errors.New("test")
		var target *CustomError
		if errors.As(err, &target) {
			t.Errorf("expected false, got true")
		}
	})
}

func TestWrap(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"only wrap": {
			err:  errors.Wrap(stderr.New("base")),
			want: "base",
		},
		"with message": {
			err:  errors.Wrap(stderr.New("base"), errors.WithMessage("wrap")),
			want: `wrap: base`,
		},
		"with no message": {
			err:  errors.Wrap(stderr.New(""), errors.WithMessage("wrap")),
			want: `wrap`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, tt.err.Error())
			}
		})
	}
}

func TestWithStack(t *testing.T) {
	t.Run("new w/ stack", func(t *testing.T) {
		err := errors.New("stack")
		var cErr *errors.Error
		if ok := errors.As(err, &cErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		if len(cErr.StackTrace()) == 0 {
			t.Errorf("expected stack trace")
		}
	})

	t.Run("wrap stderr w/ stack", func(t *testing.T) {
		baseErr := stderr.New("base")
		err := errors.Wrap(baseErr, errors.WithStack())
		var cErr *errors.Error
		if ok := errors.As(err, &cErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		if len(cErr.StackTrace()) == 0 {
			t.Errorf("expected stack trace")
		}
	})

	t.Run("take over", func(t *testing.T) {
		baseErr := errors.New("base")
		err := errors.Wrap(baseErr)
		var baseCErr, wrappedCErr *errors.Error
		if ok := errors.As(baseErr, &baseCErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		if ok := errors.As(err, &wrappedCErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		for i, pc := range baseCErr.StackTrace() {
			if wrappedCErr.StackTrace()[i] != pc {
				t.Errorf("expected same stack trace")
			}
		}
	})
}

func TestWithAttrs(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		err := errors.New("test")
		var cErr *errors.Error
		if ok := errors.As(err, &cErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		if len(cErr.Attributes()) != 0 {
			t.Errorf("expected 0 attributes, got %d", len(cErr.Attributes()))
		}
	})

	t.Run("take over", func(t *testing.T) {
		baseErr := errors.Wrap(
			errors.New("base"),
			errors.WithAttrs(
				errors.Attr("key", "value"),
			),
		)
		err := errors.Wrap(baseErr)
		var cErr *errors.Error
		if ok := errors.As(err, &cErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		if len(cErr.Attributes()) != 1 {
			t.Errorf("expected 1 attribute, got %d", len(cErr.Attributes()))
		}
		if cErr.Attributes()["key"] != "value" {
			t.Errorf("expected value, got %s", cErr.Attributes()["key"])
		}
	})

	t.Run("overwrite", func(t *testing.T) {
		baseErr := errors.Wrap(
			errors.New("base"),
			errors.WithAttrs(
				errors.Attr("key1", "value1"),
				errors.Attr("key2", "value2"),
			),
		)
		err := errors.Wrap(
			baseErr,
			errors.WithAttrs(
				errors.Attr("key1", "value111"),
				errors.Attr("key3", "value3"),
			),
		)
		var cErr *errors.Error
		if ok := errors.As(err, &cErr); !ok {
			t.Fatal("expected error to be *errors.Error")
		}
		if len(cErr.Attributes()) != 3 {
			t.Errorf("expected 1 attribute, got %d", len(cErr.Attributes()))
		}
		if cErr.Attributes()["key1"] != "value111" {
			t.Errorf("expected value, got %s", cErr.Attributes()["key1"])
		}
		if cErr.Attributes()["key2"] != "value2" {
			t.Errorf("expected value, got %s", cErr.Attributes()["key2"])
		}
		if cErr.Attributes()["key3"] != "value3" {
			t.Errorf("expected value, got %s", cErr.Attributes()["key3"])
		}
	})
}

func TestUnwrap(t *testing.T) {
	tests := map[string]struct {
		err         error
		wantNil     bool
		wantMessage string
	}{
		"single": {
			err:         errors.Wrap(errors.New("base"), errors.WithMessage("wrap")),
			wantMessage: "base",
		},
		"double": {
			err: errors.Wrap(
				errors.Wrap(
					errors.New("base"),
					errors.WithMessage("wrap"),
				),
				errors.WithMessage("wrap2"),
			),
			wantMessage: `wrap: base`,
		},
		"no wrap": {
			err:     errors.New("base"),
			wantNil: true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			err := errors.Unwrap(tt.err)
			if tt.wantNil {
				if err != nil {
					t.Errorf("expected nil, got %s", err.Error())
				}
				return
			}
			if err.Error() != tt.wantMessage {
				t.Errorf("expected %s, got %s", tt.wantMessage, err.Error())
			}
		})
	}
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
