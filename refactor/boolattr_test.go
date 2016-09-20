package refactor

import "testing"

func testBoolAttrValidateError(t *testing.T, err error, path string) {
	want := ErrInvalid(path + " must be a boolean or variable expression")
	if err == nil {
		t.Errorf("testBoolAttrValidateError() err = %v WANT %v", err, want)
		return
	}
	if err.Error() != want.Error() {
		t.Errorf("testBoolAttrValidateError() err = %v WANT %v", err, want)
	}
}
