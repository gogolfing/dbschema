package refactor

import (
	"reflect"
	"testing"

	"github.com/gogolfing/dbschema/vars"
)

func TestCreatTable_Validate_errorOnEmptyName(t *testing.T) {
	c := &CreateTable{}
	if err := c.Validate(); err != ErrInvalid("CreateTable.name cannot be empty") {
		t.Fail()
	}
}

func TestCreateTable_Validate_errorIfNotExistsNotBool(t *testing.T) {
	c := &CreateTable{Name: "name", IfNotExists: newString("")}
	err := c.Validate()
	testValidateBoolError(t, "CreateTable.ifNotExists", err)
}

func TestCreateTable_Validate_errorColumnsValidation(t *testing.T) {
	c := &CreateTable{
		Name: "name",
		Columns: []*Column{
			&Column{},
		},
	}
	if err := c.Validate(); err != ErrInvalid("Column.name cannot be empty") {
		t.Fail()
	}
}

func TestCreateTable_Up_errorValidate(t *testing.T) {
	c := &CreateTable{}
	if err := c.Validate(); err == nil {
		t.Fail()
	}
}

func TestCreateTable_Up_errorExpandAll(t *testing.T) {
	c := &CreateTable{Name: "{name}"}
	_, err := c.Up(defaultTestContext)
	if _, ok := err.(vars.ErrDoesNotExist); !ok {
		t.Fail()
	}
}

func TestCreateTable_Up_errorColumnDefinitions(t *testing.T) {
	c := &CreateTable{
		Name: "name",
		Columns: []*Column{
			&Column{},
		},
	}
	if err := c.Validate(); err != ErrInvalid("Column.name cannot be empty") {
		t.Fail()
	}
}

func TestCreateTable_Up_success(t *testing.T) {
	tests := []struct {
		c      *CreateTable
		result []Stmt
	}{
		//table without any columns.
		{
			&CreateTable{
				Name:    "table_name",
				Columns: nil,
			},
			[]Stmt{
				`CREATE TABLE "table_name" (

)`,
			},
		},

		//table with a single column.
		{
			&CreateTable{
				Name: "table_name",
				Columns: []*Column{
					&Column{Name: "col1", Type: "type1"},
				},
			},
			[]Stmt{
				`CREATE TABLE "table_name" (
	"col1" type1
)`,
			},
		},

		//table with multiple columns.
		{
			&CreateTable{
				Name: "table_name",
				Columns: []*Column{
					&Column{Name: "col1", Type: "type1"},
					&Column{Name: "col2", Type: "type2"},
				},
			},
			[]Stmt{
				`CREATE TABLE "table_name" (
	"col1" type1,
	"col2" type2
)`,
			},
		},
	}
	for _, test := range tests {
		result, err := test.c.Up(defaultTestContext)
		if err != nil {
			t.Errorf("test.c.Up(defaultTestContext) error = %v WANT %v", err, nil)
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("test.c.Up(defaultTestContext) result = %v WANT %v", result, test.result)
		}
	}
}
