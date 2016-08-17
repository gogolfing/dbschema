package dialect

import (
	"reflect"
	"testing"
)

func TestDialect_ValueOfField(t *testing.T) {
	tests := []struct {
		d      *Dialect
		field  string
		result string
		err    error
	}{
		{NewSqlDialect(), "CreateTable", "CREATE TABLE", nil},
		{NewSqlDialect(), "NotAField", "", ErrFieldDoesNotExist},
	}
	for _, test := range tests {
		result, err := test.d.ValueOfField(test.field)
		if result != test.result || !reflect.DeepEqual(err, test.err) {
			t.Fail()
		}
	}
}
