package user

import (
	"testing"
)

func TestArgon2Hash(t *testing.T) {
	hash, err := Argon2Hash("pa$$word", NewDefaultArgon2Config())
	if err != nil {
		return
	}

	invalid, err := CompareArgon2("password", hash)
	if err != nil {
		t.Fatal(err)
	}

	if invalid != false {
		t.Errorf("Password comparison returned 'true' expected 'false'")
	}

	valid, err := CompareArgon2("pa$$word", hash)
	if err != nil {
		t.Fatal(err)
	}

	if valid != true {
		t.Errorf("Password comparison returned 'false' expected 'true'")
	}
}
