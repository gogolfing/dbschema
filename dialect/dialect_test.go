package dialect

import (
	"reflect"
	"testing"
)

func TestDialect_ValueOfVariableField(t *testing.T) {
	tests := []struct {
		d       *Dialect
		varname string
		result  string
		err     error
	}{
		{NewSqlDialect(), "CreateTable", "", ErrFieldDoesNotExist},
		{NewSqlDialect(), "NotAField", "", ErrFieldDoesNotExist},

		{&Dialect{Int: "IntValue"}, "Int", "IntValue", nil},
	}
	for _, test := range tests {
		result, err := test.d.ValueOfVariableField(test.varname)
		if result != test.result || !reflect.DeepEqual(err, test.err) {
			t.Fail()
		}
	}
}
