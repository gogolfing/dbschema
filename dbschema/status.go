package dbschema

import (
	"fmt"

	"github.com/gogolfing/dbschema/logger"
	"github.com/gogolfing/dbschema/refactor"
)

func (d *DBSchema) Status(logger logger.Logger) error {
	collectingAppliedChangeSets(logger.Verbose())

	changeSetRows, err := d.listOrderedChangeSetRows()
	if err != nil {
		return err
	}

	if len(changeSetRows) == 0 {
		fmt.Fprintln(logger.Info(), "No applied ChangeSets.")
		return nil
	}

	err = iterateChangeSetRows(d.changeLog, changeSetRows, func(before []*refactor.ChangeSet, csr *ChangeSetRow) error {
		if len(before) > 0 {
			return &ErrChangeSetOutOfOrder{
				ChangeSets:   before,
				ChangeSetRow: csr,
				BeforeRow:    true,
			}
		}

		printChangeSetRow(logger, csr)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
