package refactor

import (
	"reflect"
	"testing"

	"github.com/gogolfing/dbschema/dialect"
)

func TestErrVariableDoesNotExist_Error(t *testing.T) {
	err := ErrVariableDoesNotExist("variable")
	if err.Error() != `refactor: variable does not exist "variable"` {
		t.Fail()
	}
}

func TestErrInvalidVariableReference_Error(t *testing.T) {
	err := ErrInvalidVariableReference("reference")
	if err.Error() != `refactor: invalid variable reference "reference"` {
		t.Fail()
	}
}

func TestNewContext(t *testing.T) {
	d := &dialect.Dialect{}
	ctx := NewContext(d)

	if ctx.Dialect != d || ctx.Variables == nil || len(ctx.Variables) != 0 {
		t.Fail()
	}
}

func TestContext_GetVariable(t *testing.T) {
	tests := []struct {
		c      *Context
		name   string
		result string
		err    error
	}{
		{
			NewContext(nil),
			"{Dialect.UUID}",
			"UUID",
			nil,
		},
		{
			func() *Context {
				c := NewContext(nil)
				c.Variables["foo"] = "bar"
				return c
			}(),
			"{foo}",
			"bar",
			nil,
		},
		{
			func() *Context {
				c := NewContext(nil)
				c.Variables["foo"] = "bar"
				return c
			}(),
			"foo",
			"",
			ErrInvalidVariableReference("foo"),
		},
		{
			NewContext(nil),
			"{Dialect.DoesNotExist}",
			"",
			ErrVariableDoesNotExist("Dialect.DoesNotExist"),
		},
		{
			NewContext(nil),
			"{foo}",
			"",
			ErrVariableDoesNotExist("foo"),
		},
	}
	for _, test := range tests {
		result, err := test.c.GetVariable(test.name)
		if result != test.result || !reflect.DeepEqual(err, test.err) {
			t.Errorf("test.c.GetVariable(%v) = %v, %v WANT %v, %v", test.name, result, err, test.result, test.err)
		}
	}
}
