package dbschema

import (
	"database/sql"

	"github.com/gogolfing/dbschema/refactor"
)

const (
	DefaultChangeLogTableName     = "dbschema_changelog"
	DefaultChangeLogLockTableName = "dbschema_changelog_lock"
)

func (d *DBSchema) init() error {
	if err := d.initTables(); err != nil {
		return err
	}
	if err := d.initLock(); err != nil {
		return err
	}
	return nil
}

func (d *DBSchema) initTables() error {
	expanded, err := refactor.ExpandAll(
		d,
		d.changeLog.ChangeLogTableName.ExpandDefault(DefaultChangeLogTableName),
		d.changeLog.ChangeLogLockTableName.ExpandDefault(DefaultChangeLogLockTableName),
	)
	if err != nil {
		return err
	}

	changeLogTableName, changeLogLockTableName := expanded[0], expanded[1]

	changeLogTable := createChangeLogTable(changeLogTableName)
	changeLogLockTable := createChangeLogLockTable(changeLogLockTableName)

	return d.executeTxChangers(changeLogTable, changeLogLockTable)
}

func (d *DBSchema) initLock() error {
	work := func(tx *sql.Tx) error {
		return nil
	}
	return d.executeTxWork(work)
}
