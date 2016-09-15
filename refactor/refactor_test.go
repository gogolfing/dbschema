package refactor

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gogolfing/dbschema/refactor/strval"
)

func TestErrInvalid_Error(t *testing.T) {
	err := ErrInvalid("this is a validation error")

	if err.Error() != "dbschema/refactor: invalid: this is a validation error" {
		t.Fail()
	}
}

func decodeSourceIntoValue(t *testing.T, value interface{}, source string) {
	dec := xml.NewDecoder(strings.NewReader(source))
	err := dec.Decode(value)
	if err != nil {
		t.Fatal(err)
	}
}

func decodeSourceIntoValueError(t *testing.T, value interface{}, source string) error {
	dec := xml.NewDecoder(strings.NewReader(source))
	err := dec.Decode(value)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	return err
}

func newString(s string) *string {
	return &s
}

func writeFile(t *testing.T, file *os.File, source string) {
	if _, err := file.WriteString(source); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func testValidateBoolError(t *testing.T, prefix string, err error) {
	if err == nil {
		t.Error("err must not be nil")
	}
	boolErr := strval.ValidateBool(newString(""))
	want := fmt.Errorf("%v %v", prefix, boolErr)
	if want.Error() != err.Error() {
		t.Errorf("err = %v WANT %v", err, want)
	}
}
