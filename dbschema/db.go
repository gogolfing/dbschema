package dbschema

import (
	"database/sql"
	"fmt"
	"io"
)

type UnsupportedDialectError string

func (e UnsupportedDialectError) Error() string {
	return fmt.Sprintf("dbschema/dbschema: unsupported dialect.Dialect.DBMS() %q", string(e))
}

type DB interface {
	Ping() error
	Begin() (Tx, error)
	io.Closer
}

type Tx interface {
	QueryExecer
	Rollback() error
	Commit() error
}

type QueryExecer interface {
	Queryer
	Execer
}

type Queryer interface {
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Row interface {
	Scanner
}

type Rows interface {
	io.Closer
	Next() bool
	Columns() ([]string, error)
	Err() error
	Scanner
}

type Scanner interface {
	Scan(dest ...interface{}) error
}
