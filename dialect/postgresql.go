package dialect

import "github.com/gogolfing/dbschema/conn"

const Postgresql = "postgresql"

func NewDialectPostgresql() Dialect {
	return newDialect(newPostgresqlConnectionString, newPostgresqlDialectStruct())
}

func newPostgresqlConnectionString(conn *conn.Connection) (string, error) {
	return "", nil
}

func newPostgresqlDialectStruct() *DialectStruct {
	return &DialectStruct{
		DBMSValue: Postgresql,

		QuoteRefValue:   `"`,
		QuoteConstValue: "'",

		IntValue: "INTEGER",

		UUIDValue: "UUID",
	}
}
