package dbschema

import (
	"fmt"

	"github.com/gogolfing/dbschema/logger"
	"github.com/gogolfing/dbschema/refactor"
)

func (d *DBSchema) Status(logger logger.Logger) error {
	fmt.Fprintln(logger.Verbose(), "collecting applied ChangeSets...")

	changeSetRows, err := d.collectChangeSetRows()
	if err != nil {
		return err
	}

	for _, csr := range changeSetRows {
		fmt.Fprintln(logger.Info(), csr)
	}

	return nil
}

func (d *DBSchema) collectChangeSets() ([]*refactor.ChangeSet, error) {
	return nil, nil
}
