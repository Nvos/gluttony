package passwordutil_test

import (
	"gluttony/internal/util/passwordutil"
	"testing"
)

func TestArgon2Hash(t *testing.T) {
	hash, err := passwordutil.Argon2Hash("pa$$word", passwordutil.NewDefaultArgon2Config())
	if err != nil {
		return
	}

	invalid, err := passwordutil.CompareArgon2("password", hash)
	if err != nil {
		t.Fatal(err)
	}

	if invalid != false {
		t.Errorf("Password comparison returned 'true' expected 'false'")
	}

	valid, err := passwordutil.CompareArgon2("pa$$word", hash)
	if err != nil {
		t.Fatal(err)
	}

	if valid != true {
		t.Errorf("Password comparison returned 'false' expected 'true'")
	}
}
