package mysql

import (
	"github.com/gogolfing/dbschema/src/dialect"
)

const DBMS = "mysql"

const DefaultPort = 3306

func Dialect() dialect.Dialect {
	return &dialect.DialectStruct{
		DBMSValue: DBMS,
		Syntax:    Syntax(),
		Types:     Types(),
	}
}

func Syntax() dialect.Syntax {
	return &dialect.SyntaxStruct{
		QuoteRefValue:   "`",
		QuoteConstValue: `'`,

		Escapes: newEscapes(),

		//TODO change.
		Caster: dialect.DoubleColonCaster,

		//TODO change.
		PlaceholderValue: func(num int) string { return "?" },
	}
}

func Types() dialect.Types {
	return &dialect.TypesStruct{
		IntegerValue: "INTEGER",
		Int8Value:    "TINYINT",
		Int16Value:   "SMALLINT",
		Int32Value:   "INTEGER",
		Int64Value:   "BIGINT",

		Float32Value: "FLOAT",
		Float64Value: "DOUBLE PRECISION",

		Char32Value:  "CHAR(32)",
		Char64Value:  "CHAR(64)",
		Char128Value: "CHAR(128)",
		Char256Value: "CHAR(256)",

		VarChar32Value:   "VARCHAR(32)",
		VarChar64Value:   "VARCHAR(64)",
		VarChar128Value:  "VARCHAR(128)",
		VarChar256Value:  "VARCHAR(256)",
		VarChar512Value:  "VARCHAR(512)",
		VarChar1024Value: "VARCHAR(1024)",

		TextValue: "TEXT",

		ByteArrayValue: "BLOB",

		TimestampValue:   "TIMESTAMP",
		TimestampTzValue: "TIMESTAMP",
		TimeValue:        "TIME",
		TimeTzValue:      "TIME",
		DateValue:        "DATE",

		BoolValue: "BOOLEAN",

		UUIDValue: "VARCHAR(36)",
	}
}

func newEscapes() map[string]string {
	escapes := dialect.NewDefaultEscapes()
	escapes["`"] = "\\`"
	return escapes
}
