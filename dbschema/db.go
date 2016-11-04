package dbschema

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/refactor"
)

type ErrUnsupportedDialect string

func (e ErrUnsupportedDialect) Error() string {
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
	Querier
	Execer
}

type Querier interface {
	Query(stmt *refactor.Stmt) (Rows, error)
	QueryRow(stmt *refactor.Stmt) Row
}

type Execer interface {
	Exec(stmt *refactor.Stmt) (sql.Result, error)
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
