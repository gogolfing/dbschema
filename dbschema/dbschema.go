package dbschema

import (
	"database/sql"
	"fmt"

	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/refactor"

	_ "github.com/lib/pq"
)

type DbSchema struct {
	dialect.Dialect

	changeLog *refactor.ChangeLog

	refactorContext refactor.Context

	db DB
}

func OpenDbSchemaSql(dialect dialect.Dialect, db *sql.DB, changeLog *refactor.ChangeLog) (*DbSchema, error) {
	d := &DbSchema{
		Dialect:   dialect,
		changeLog: changeLog,
	}
	if err := d.openSql(db); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DbSchema) openSql(db *sql.DB) error {
	var err error = nil
	if err = db.Ping(); err != nil {
		return err
	}
	d.db, err = createDB(db)
	return err
}

func createDB(db *sql.DB) (DB, error) {
	return nil, nil
}

func (d *DbSchema) Close() error {
	return fmt.Errorf("unimplemented")
}

func (d *DbSchema) Up(log Logger, count int) error {
	return nil
}

func (d *DbSchema) Expand(expr string) (value string, err error) {
	value, err = d.changeLog.Variables.Dereference(expr)
	return
}
