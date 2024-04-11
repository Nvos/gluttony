package cryptox_test

import (
	"gluttony/internal/x/cryptox"
	"testing"
)

func TestArgon2Hash(t *testing.T) {
	hash, err := cryptox.Argon2Hash("pa$$word", cryptox.NewDefaultArgon2Config())
	if err != nil {
		return
	}

	invalid, err := cryptox.CompareArgon2("password", hash)
	if err != nil {
		t.Fatal(err)
	}

	if invalid != false {
		t.Errorf("Password comparison returned 'true' expected 'false'")
	}

	valid, err := cryptox.CompareArgon2("pa$$word", hash)
	if err != nil {
		t.Fatal(err)
	}

	if valid != true {
		t.Errorf("Password comparison returned 'false' expected 'true'")
	}
}
