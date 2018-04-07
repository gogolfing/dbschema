package dbschema

import (
	"database/sql"

	"github.com/gogolfing/dbschema/src/dialect"
	"github.com/gogolfing/dbschema/src/dialect/postgresql"
)

var (
	DriverPostgresql = "postgresql"
)

type sqlDB struct {
	db *sql.DB
}

func openSqlDB(d dialect.Dialect, conn string) (DB, error) {
	dbmsToDriver := map[string]string{
		postgresql.DBMS: DriverPostgresql,
	}
	driverName, ok := dbmsToDriver[d.DBMS()]
	if !ok {
		return nil, UnsupportedDialectError(d.DBMS())
	}
	db, err := sql.Open(driverName, conn)
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
	return &sqlTx{tx}, nil
}

func (s *sqlDB) Close() error {
	return s.db.Close()
}

type sqlTx struct {
	*sql.Tx
}

func (s *sqlTx) Query(query string, args ...interface{}) (Rows, error) {
	return s.Tx.Query(query, args...)
}

func (s *sqlTx) QueryRow(query string, args ...interface{}) Row {
	return s.Tx.QueryRow(query, args...)
}
