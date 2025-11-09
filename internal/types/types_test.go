package types

import "testing"

func TestErrorWrapper(t *testing.T) {
	empty := Error{}
	if empty.Error() != "" {
		t.Fatalf("expected empty error string")
	}

	err := Error{Code: ExitCodeData, Err: ErrExitCode{}}
	if err.Error() != "data error" {
		t.Fatalf("unexpected error string: %s", err.Error())
	}
	if err.Unwrap() == nil {
		t.Fatalf("expected unwrap to return underlying error")
	}
}

type ErrExitCode struct{}

func (ErrExitCode) Error() string { return "data error" }
