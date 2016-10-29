package dbschema

import (
	"fmt"

	"github.com/gogolfing/dbschema/refactor"
)

const DefaultChangeLogTableName = "dbschema_changelog"

func (d *DBSchema) init() error {
	return d.initTables()
}

func (d *DBSchema) initTables() error {
	expanded, err := refactor.ExpandAll(
		d,
		d.changeLog.ChangeLogTableName.ExpandDefault(DefaultChangeLogTableName),
	)
	if err != nil {
		return err
	}

	changeLogTableName := expanded[0]

	fmt.Println(d.changeLog)

	changeLogTable := createChangeLogTable(changeLogTableName)

	return d.executeInTransaction(changeLogTable)
}
