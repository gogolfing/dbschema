package dbschema

import (
	"database/sql"
	"fmt"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/refactor"
)

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

func (s *sqlDB) Ping() error {
	return s.db.Ping()
}

func (s *sqlDB) Begin() (Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return SqlTx{tx}, nil
}

func (s *sqlDB) Close() error {
	return s.db.Close()
}

type SqlTx struct {
	*sql.Tx
}

func (s SqlTx) Exec(stmt refactor.Stmt) (sql.Result, error) {
	fmt.Println(stmt)
	return s.Tx.Exec(string(stmt))
}

func (s SqlTx) Query(stmt refactor.Stmt) (Rows, error) {
	return s.Tx.Query(string(stmt))
}

func (s SqlTx) QueryRow(stmt refactor.Stmt) Row {
	return s.Tx.QueryRow(string(stmt))
}

func (s SqlTx) Commit() error {
	return s.Tx.Commit()
}

func (s SqlTx) Rollback() error {
	return s.Tx.Rollback()
}
