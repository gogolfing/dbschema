package refactor

import (
	"reflect"
	"testing"
)

func TestCreateTable_Up_error(t *testing.T) {
	a := &CreateTable{}

	result, err := a.Up(NewContext(nil))

	if _, ok := err.(ErrInvalid); result != nil || !ok {
		t.Fail()
	}
}

func TestCreateTable_Up_emptyColumns(t *testing.T) {
	a := &CreateTable{
		Name: "TableName",
	}

	result, err := a.Up(NewContext(nil))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(
		result,
		[]string{
			"CREATE TABLE TableName ()",
		},
	) {
		t.Fail()
	}
}
