package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) bool {
	t.Helper()
	if actual != expected {
		t.Errorf("got %v, want: %v", actual, expected)
		return true
	}

	return false
}

func NilErr(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("got %v error, want nil", err)
		return true
	}

	return false
}

func Err(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		t.Errorf("got nil error")
		return true
	}

	return false
}
