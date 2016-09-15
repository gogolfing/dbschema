package refactor

import (
	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
	_vars "github.com/gogolfing/dbschema/vars"
)

const (
	DefaultTestContextDBMS             = "dbscema/refactor/context_test.go/dbms"
	DefaultTestContextConnectionString = "dbscema/refactor/context_test.go/connection_string"
)

var defaultTestDialect *testDialect

var defaultTestContext Context

func init() {
	defaultTestDialect = &testDialect{
		DialectStruct: newDefaultDialectStruct(),
	}
	defaultTestContext = newTestContext(defaultTestDialect, []string{}...)
}

func newTestContext(d dialect.Dialect, vars ...string) Context {
	if len(vars)%2 == 1 {
		vars = vars[:len(vars)-1]
	}
	v := &_vars.Variables{}
	for i := 0; i < len(vars); i += 2 {
		v.Put(&_vars.Variable{Name: vars[i], Value: vars[i+1]})
	}
	return &context{Dialect: d, vars: v}
}

type context struct {
	dialect.Dialect

	vars *_vars.Variables
}

func (c *context) Expand(expr string) (value string, err error) {
	return dialect.Expand(expr, c.vars, c.Dialect)
}

func newDefaultTestDialect() dialect.Dialect {
	return &testDialect{
		DialectStruct: newDefaultDialectStruct(),
	}
}

type testDialect struct {
	*dialect.DialectStruct
}

func (td *testDialect) ConnectionString(_ *conn.Connection) (string, error) {
	return DefaultTestContextConnectionString, nil
}

func newDefaultDialectStruct() *dialect.DialectStruct {
	return &dialect.DialectStruct{
		DBMSValue: DefaultTestContextDBMS,

		QuoteRefValue:   `"`,
		QuoteConstValue: "'",

		IntegerValue: "test_integer",
		Int8Value:    "test_int8",
		Int16Value:   "test_int16",
		Int32Value:   "test_int32",
		Int64Value:   "test_int64",

		Float32Value: "test_float32",
		Float64Value: "test_float64",

		Char64Value:  "test_char64",
		Char128Value: "test_char128",
		Char256Value: "test_char256",

		VarChar64Value:  "test_varchar64",
		VarChar128Value: "test_varchar128",
		VarChar256Value: "test_varchar256",

		TextValue: "test_text",

		TimestampValue:   "test_timestamp",
		TimestampTzValue: "test_timestamp_with_time_zone",
		TimeValue:        "test_time",
		TimeTzValue:      "test_time_with_time_zone",
		DateValue:        "test_date",

		BoolValue: "test_bool",

		UUIDValue: "test_uuid",
	}
}
