package refactor

import (
	"testing"

	"github.com/gogolfing/dbschema/vars"
)

func NewColumn(name, t string) *Column {
	return &Column{
		Name: NewStringAttr(name),
		Type: NewStringAttr(t),
	}
}

func TestColumn_Validate_errorOnEmtpyName(t *testing.T) {
	c := &Column{}
	if err := c.Validate(); err != ErrInvalid("Column.name cannot be empty") {
		t.Fail()
	}
}

func TestColumn_Validate_errorOnEmtpyType(t *testing.T) {
	c := &Column{Name: NewStringAttr("name")}
	if err := c.Validate(); err != ErrInvalid("Column.type cannot be empty") {
		t.Fail()
	}
}

func TestColumn_Validate_errorOnIsNullableNotBool(t *testing.T) {
	c := &Column{
		Name:       NewStringAttr("name"),
		Type:       NewStringAttr("type"),
		IsNullable: NewBoolAttr("foo"),
	}
	err := c.Validate()
	testBoolAttrValidateError(t, err, "Column.isNullable")
}

func TestColumn_Validate_errorOnConstraintValidate(t *testing.T) {
	c := NewColumn("name", "type")
	c.Constraint = invalidConstraint
	if err := c.Validate(); err != invalidConstraintError {
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
	c := &Column{Name: NewStringAttr("${name}"), Type: NewStringAttr("type")}
	_, err := c.Definition(defaultTestContext)
	if _, ok := err.(vars.ErrDoesNotExist); !ok {
		t.Fail()
	}
}

func TestColumn_Definition_success(t *testing.T) {
	tests := []struct {
		ctx    Context
		c      *Column
		result string
	}{
		{
			defaultTestContext,
			&Column{
				Name: NewStringAttr("name"),
				Type: NewStringAttr("type"),
			},
			`"name" type`,
		},
		{
			newTestContext(defaultTestDialect, "var_name", "var_value"),
			&Column{
				Name: NewStringAttr("${var_name}"),
				Type: NewStringAttr("type"),
			},
			`"var_value" type`,
		},
		{
			newTestContext(defaultTestDialect, "var_name", "var_value"),
			&Column{
				Name:    NewStringAttr("name"),
				Type:    NewStringAttr("type"),
				Default: NewStringAttr("${var_name}"),
			},
			`"name" type DEFAULT 'var_value'::type`,
		},
		{
			defaultTestContext,
			&Column{
				Name:       NewStringAttr("name"),
				Type:       NewStringAttr("${Dialect.Bool}"),
				IsNullable: NewBoolAttr("false"),
			},
			`"name" ` + defaultTestDialect.BoolValue + ` NOT NULL`,
		},
		{
			defaultTestContext,
			&Column{
				Name:    NewStringAttr("name"),
				Type:    NewStringAttr("type"),
				Default: NewStringAttr("default_value"),
			},
			`"name" type DEFAULT 'default_value'::type`,
		},
		{
			defaultTestContext,
			&Column{
				Name:       NewStringAttr("name"),
				Type:       NewStringAttr("type"),
				IsNullable: NewBoolAttr("false"),
				Default:    NewStringAttr("def'ault\nvalue"),
			},
			`"name" type NOT NULL DEFAULT 'def\'ault\nvalue'::type`,
		},
	}
	for _, test := range tests {
		result, err := test.c.Definition(test.ctx)
		if err != nil {
			t.Errorf("test.c.Definition(test.ctx) error = %v WANT %v", err, nil)
		}
		if result != test.result {
			t.Errorf("test.c.Definition(test.ctx) result = %v WANT %v", result, test.result)
		}
	}
}
