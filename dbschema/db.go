package dbschema

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
)

type ErrUnsupportedDialect string

func (e ErrUnsupportedDialect) Error() string {
	return fmt.Sprintf("dbschema/dbchema: unsupported dialect.Dialect.DBMS() %q", string(e))
}

type DB interface {
	Begin() (*sql.Tx, error)
	Ping() error
	Exec(query string, args ...interface{}) (sql.Result, error)

	QueryDB

	io.Closer
}

type QueryDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type sqlDB struct {
	db *sql.DB
}

func openSqlDB(d dialect.Dialect, conn *conn.Connection) (DB, error) {
	dbmsToDriver := map[string]string{
		dialect.Postgresql: "postgres",
	}
	driverName, ok := dbmsToDriver[d.DBMS()]
	if !ok {
		return nil, ErrUnsupportedDialect(d.DBMS())
	}
	connString, err := d.ConnectionString(conn)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return nil, err
	}
	return &sqlDB{db}, nil
}

func (s *sqlDB) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *sqlDB) Ping() error {
	return s.db.Ping()
}

func (s *sqlDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *sqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *sqlDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *sqlDB) Close() error {
	return s.db.Close()
}
