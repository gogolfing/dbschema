package refactor

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
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

/*
func testUpDownResult(t *testing.T, cf func(Context) ([]string, error), ctx Context, err error, stmts ...string) {
	result, resultErr := cf(ctx)
	if !reflect.DeepEqual(result, stmts) || !reflect.DeepEqual(resultErr, err) {
		t.Errorf("cf(ctx) = %v, %v WANT %v, %v", result, resultErr, stmts, err)
	}
}
*/
