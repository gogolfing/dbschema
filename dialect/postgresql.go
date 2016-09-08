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
	url := &url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   host,
		Path:   path,
	}
	return url.String(), nil
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
