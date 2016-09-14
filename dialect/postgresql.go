package dialect

import (
	"net/url"

	"github.com/gogolfing/dbschema/conn"
)

const Postgresql = "postgresql"

const PostgresqlDefaultPort = 5432

func NewDialectPostgresql() Dialect {
	return newDialect(newPostgresqlConnectionString, newPostgresqlDialectStruct())
}

func newPostgresqlConnectionString(conn *conn.Connection) (string, error) {
	user, err := conn.Userinfo()
	if err != nil {
		return "", err
	}
	host, err := conn.HostPort(PostgresqlDefaultPort)
	if err != nil {
		return "", err
	}
	path, err := conn.DatabaseValue()
	if err != nil {
		return "", err
	}
	rawQuery, err := conn.Query()
	if err != nil {
		return "", err
	}
	url := &url.URL{
		Scheme:   "postgres",
		User:     user,
		Host:     host,
		Path:     path,
		RawQuery: rawQuery,
	}
	return url.String(), nil
}

func newPostgresqlDialectStruct() *DialectStruct {
	return &DialectStruct{
		DBMSValue: Postgresql,

		QuoteRefValue:   `"`,
		QuoteConstValue: "'",

		IntegerValue: "INTEGER",
		// Int8Value not implemented.
		Int16Value: "SMALLINT",
		Int32Value: "INTEGER",
		Int64Value: "BIGINT",

		Float32Value: "REAL",
		Float64Value: "DOUBLE PRECISION",

		Char64Value:  "CHARACTER(64)",
		Char128Value: "CHARACTER(128)",
		Char256Value: "CHARACTER(256)",

		VarChar64Value:  "CHARACTER VARYING(64)",
		VarChar128Value: "CHARACTER VARYING(128)",
		VarChar256Value: "CHARACTER VARYING(256)",

		TextValue: "TEXT",

		TimestampValue:   "TIMESTAMP",
		TimestampTzValue: "TIMESTAMP WITH TIME ZONE",
		TimeValue:        "TIME",
		TimeTzValue:      "TIME WITH TIME ZONE",
		DateValue:        "DATE",

		BoolValue: "BOOLEAN",

		UUIDValue: "UUID",
	}
}
