package dialect

import "github.com/gogolfing/dbschema/conn"

const Postgresql = "postgresql"

func NewPostgresqlDialect() Dialect {
	return newDialect(newPostgresqlConnectionString, newPostgresqlDialectStruct())
}

func newPostgresqlConnectionString(conn *conn.Connection) (string, error) {
	return "", nil
}

func newPostgresqlDialectStruct() *DialectStruct {
	return &DialectStruct{
		QuoteRefValue: `"`,

		QuoteConstValue: "'",

		IntValue: "INTEGER",

		UUIDValue: "UUID",
	}
}
