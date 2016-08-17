package refactor

import (
	"reflect"
	"testing"
)

func TestAddTable_Up_error(t *testing.T) {
	a := &AddTable{}

	result, err := a.Up(NewContext(nil))

	if _, ok := err.(ValidationError); result != nil || !ok {
		t.Fail()
	}
}

func TestAddTable_Up_emptyColumns(t *testing.T) {
	a := &AddTable{
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
