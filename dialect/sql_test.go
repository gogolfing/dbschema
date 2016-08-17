package dialect

import (
	"reflect"
	"testing"
)

func TestNewSqlDialect(t *testing.T) {
	d := NewSqlDialect()
	if !reflect.DeepEqual(
		d,
		&Dialect{
			CreateTable: "CREATE TABLE",

			Int:  "int",
			UUID: "UUID",
		},
	) {
		t.Fail()
	}
}
