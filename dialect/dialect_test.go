package dialect

import (
	"reflect"
	"testing"
)

func TestDialectStruct_isDialect(t *testing.T) {
	var d Dialect = nil
	d = &DialectStruct{}
	d.Int() //to avoid the not used error.
}

func TestCallVariableMethodOnDialect(t *testing.T) {
	d := &DialectStruct{
		IntValue: "int",
	}

	tests := []struct {
		name      string
		result    string
		resultErr error
	}{
		{"doesNotExist", "", ErrMethodDoesNotExist},
		{"QuoteRef", "", ErrInvalidVariableMethodType},
		{"UUID", "", ErrNotSupported},
		{"Int", "int", nil},
	}
	for _, test := range tests {
		result, resultErr := CallVariableMethodOnDialect(d, test.name)
		if result != test.result || !reflect.DeepEqual(resultErr, test.resultErr) {
			t.Errorf(
				"CallVariableMethodOnDialect(%v, %v) = %v, %v WANT %v, %v",
				d,
				test.name,
				result,
				resultErr,
				test.result,
				test.resultErr,
			)
		}
	}
}
