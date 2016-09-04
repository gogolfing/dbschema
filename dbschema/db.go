package dbschema

import (
	"database/sql"
	"io"

	"github.com/gogolfing/dbschema/refactor"
)

type DB interface {
	Query(stmt refactor.Stmt, args ...interface{}) (*sql.Rows, error)
	QueryRow(stmt refactor.Stmt, args ...interface{}) (*sql.Result, error)

	io.Closer
}

type db *sql.DB

func (db db) Query(stmt refactor.Stmt, args ...interface{}) (*sql.Rows, error) {
	return db.Query(stmt.String(), args...)
}
