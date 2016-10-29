package dbschema

import (
	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/refactor"

	_ "github.com/lib/pq"
)

type DBSchema struct {
	db DB

	dialect.Dialect

	changeLog *refactor.ChangeLog
}

func OpenSql(dialect dialect.Dialect, conn *conn.Connection, changeLog *refactor.ChangeLog) (*DBSchema, error) {
	sqlDB, err := openSqlDB(dialect, conn)
	if err != nil {
		return nil, err
	}
	return openDB(sqlDB, dialect, changeLog)
}

func openDB(db DB, dialect dialect.Dialect, changeLog *refactor.ChangeLog) (*DBSchema, error) {
	d := &DBSchema{
		db:        db,
		Dialect:   dialect,
		changeLog: changeLog,
	}
	if err := d.open(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DBSchema) open() error {
	if err := d.db.Ping(); err != nil {
		return err
	}
	if err := d.init(); err != nil {
		return err
	}
	return nil
}

func (d *DBSchema) Close() error {
	if err := d.finalize(); err != nil {
		return err
	}
	return d.db.Close()
}

func (d *DBSchema) finalize() error {
	return nil
}

func (d *DBSchema) Expand(expr string) (value string, err error) {
	return dialect.Expand(expr, d.changeLog.Variables, d.Dialect)
}

func (d *DBSchema) executeInTransaction(changers ...refactor.Changer) (err error) {
	stmts, err := d.collectChangerStmts(changers...)
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, stmt := range stmts {
		_, err = tx.Exec(string(stmt))
		if err != nil {
			return
		}
	}

	err = tx.Commit()
	return
}

func (d *DBSchema) collectChangerStmts(changers ...refactor.Changer) ([]refactor.Stmt, error) {
	result := []refactor.Stmt{}
	for _, changer := range changers {
		stmts, err := changer.Stmts(d)
		if err != nil {
			return nil, err
		}
		result = append(result, stmts...)
	}
	return result, nil
}
