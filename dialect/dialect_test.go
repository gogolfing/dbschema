package dialect

import (
	"reflect"
	"testing"

	"github.com/gogolfing/dbschema/conn"
)

func TestCallVariableMethodOnDialect(t *testing.T) {
	d := newDialect(
		func(_ *conn.Connection) (string, error) { return "connection_string", nil },
		&DialectStruct{
			IntegerValue: "int",
			BoolValue:    "boolean",
		},
	)

	tests := []struct {
		name      string
		result    string
		resultErr error
	}{
		{"doesNotExist", "", ErrMethodDoesNotExist},
		{"QuoteRef", "", ErrInvalidVariableMethodType},
		{"UUID", "", ErrNotSupported},
		{"Integer", "int", nil},
		{"Bool", "boolean", nil},
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
