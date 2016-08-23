package refactor

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

func TestStringDefaultBool(t *testing.T) {
	tr, f, notTrue := "true", "false", "not true"

	tests := []struct {
		value  *string
		def    bool
		result bool
	}{
		{nil, true, true},
		{nil, false, false},

		{&notTrue, true, false},
		{&notTrue, false, false},

		{&f, true, false},
		{&f, false, false},

		{&tr, true, true},
		{&tr, false, true},
	}
	for _, test := range tests {
		result := StringDefaultBool(test.value, test.def)
		if result != test.result {
			t.Fail()
		}
	}
}

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

func testUpDownResult(t *testing.T, cf func(Context) ([]string, error), ctx Context, err error, stmts ...string) {
	result, resultErr := cf(ctx)
	if !reflect.DeepEqual(result, stmts) || !reflect.DeepEqual(resultErr, err) {
		t.Errorf("cf(ctx) = %v, %v WANT %v, %v", result, resultErr, stmts, err)
	}
}
