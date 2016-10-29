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
		QuoteConstValue: `'`,

		Escapes: newPostgresqlEscapes(),

		Caster: DoubleColonCaster,

		IntegerValue: "INTEGER",
		Int8Value:    "SMALLINT", //this is the same as the Int16Value because Int8 is not implemented. users may override elsewhere.
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

		TimestampValue:   "TIMESTAMP",
		TimestampTzValue: "TIMESTAMP WITH TIME ZONE",
		TimeValue:        "TIME",
		TimeTzValue:      "TIME WITH TIME ZONE",
		DateValue:        "DATE",

		BoolValue: "BOOLEAN",

		UUIDValue: "UUID",
	}
}

func newPostgresqlEscapes() map[string]string {
	escapes := NewDefaultEscapes()
	escapes[`'`] = `\'`
	return escapes
}
