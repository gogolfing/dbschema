package refactor

import "testing"

func testStringAttrEmptyError(t *testing.T, err error, prefix string) {
	want := ErrInvalid(prefix + " cannot be empty")
	if err == nil {
		t.Errorf("testStringAttrEmptyError() err = %v WANT %v", err, want)
		return
	}
	if err.Error() != want.Error() {
		t.Errorf("testStringAttrEmptyError() err = %v WANT %v", err, want)
	}
}
