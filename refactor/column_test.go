package refactor

import (
	"fmt"
	"testing"

	"github.com/gogolfing/dbschema/refactor/strval"
	"github.com/gogolfing/dbschema/vars"
)

func TestColumn_Validate_errorOnEmtpyName(t *testing.T) {
	c := &Column{}
	if err := c.Validate(); err != ErrInvalid("Column.name cannot be empty") {
		t.Fail()
	}
}

func TestColumn_Validate_errorOnEmtpyType(t *testing.T) {
	c := &Column{Name: "name"}
	if err := c.Validate(); err != ErrInvalid("Column.type cannot be empty") {
		t.Fail()
	}
}

func TestColumn_Validate_errorOnIsNullableNotBool(t *testing.T) {
	c := &Column{Name: "name", Type: "type", IsNullable: newString("")}
	err := c.Validate()
	testValidateBoolError(t, "Column.isNullable", err)
}

func TestColumn_Validate_errorOnConstraintValidate(t *testing.T) {
	c := &Column{Name: "name", Type: "type", Constraint: invalidConstraint}
	if err := c.Validate(); err == nil {
		t.Fail()
	}
}

func TestColumn_Definition_errorValidate(t *testing.T) {
	c := &Column{}
	if _, err := c.Definition(defaultTestContext); err == nil {
		t.Fail()
	}
}

func TestColumn_Definition_errorExpand(t *testing.T) {
	c := &Column{Name: "{name}", Type: "type"}
	_, err := c.Definition(defaultTestContext)
	if _, ok := err.(vars.ErrDoesNotExist); !ok {
		t.Fail()
	}
}

func TestColumn_Definition_success(t *testing.T) {
	tests := []struct {
		c      *Column
		result string
	}{
		{
			&Column{
				Name: "name",
				Type: "type",
			},
			"name type",
		},
		{
			&Column{
				Name:       "name",
				Type:       "{Dialect.Bool}",
				IsNullable: newString(strval.False),
			},
			fmt.Sprintf("name %v NOT NULL", defaultTestDialect.BoolValue),
		},
	}
	for _, test := range tests {
		result, err := test.c.Definition(defaultTestContext)
		if err != nil {
			t.Errorf("test.c.Definition(defaultContext) error = %v WANT %v", err, nil)
		}
		if result != test.result {
			t.Errorf("test.c.Definition(defaultContext) result = %v WANT %v", result, test.result)
		}
	}
}