package postgresql

import (
	"fmt"

	"github.com/gogolfing/dbschema/src/dialect"
)

const DBMS = "postgresql"

const DefaultPort = 5432

func Dialect() dialect.Dialect {
	return &dialect.DialectStruct{
		DBMSValue: DBMS,
		Syntax:    Syntax(),
		Types:     Types(),
	}
}

func Syntax() dialect.Syntax {
	return &dialect.SyntaxStruct{
		QuoteRefValue:   `"`,
		QuoteConstValue: `'`,

		Escapes: newEscapes(),

		Caster: dialect.DoubleColonCaster,

		PlaceholderValue: func(num int) string { return fmt.Sprintf("$%v", num+1) },
	}
}

func Types() dialect.Types {
	return &dialect.TypesStruct{
		IntegerValue: "INTEGER",
		Int8Value:    "SMALLINT", //This is the same as the Int16Value because Int8 is not implemented. Users may override elsewhere.
		Int16Value:   "SMALLINT",
		Int32Value:   "INTEGER",
		Int64Value:   "BIGINT",

		Float32Value: "REAL",
		Float64Value: "DOUBLE PRECISION",

		Char32Value:  "CHARACTER(32)",
		Char64Value:  "CHARACTER(64)",
		Char128Value: "CHARACTER(128)",
		Char256Value: "CHARACTER(256)",

		VarChar32Value:   "CHARACTER VARYING(32)",
		VarChar64Value:   "CHARACTER VARYING(64)",
		VarChar128Value:  "CHARACTER VARYING(128)",
		VarChar256Value:  "CHARACTER VARYING(256)",
		VarChar512Value:  "CHARACTER VARYING(512)",
		VarChar1024Value: "CHARACTER VARYING(1024)",

		TextValue: "TEXT",

		ByteArrayValue: "BYTEA",

		TimestampValue:   "TIMESTAMP",
		TimestampTzValue: "TIMESTAMP WITH TIME ZONE",
		TimeValue:        "TIME",
		TimeTzValue:      "TIME WITH TIME ZONE",
		DateValue:        "DATE",

		BoolValue: "BOOLEAN",

		UUIDValue: "UUID",
	}
}

func newEscapes() map[string]string {
	escapes := dialect.NewDefaultEscapes()
	escapes[`'`] = `\'`
	return escapes
}
