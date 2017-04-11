package dbschema

import (
	"fmt"

	"github.com/gogolfing/dbschema/logger"
	"github.com/gogolfing/dbschema/refactor"
)

func (d *DBSchema) Up(logger logger.Logger, count int) error {
	collectingAppliedChangeSets(logger.Verbose())

	changeSetRows, err := d.listOrderedChangeSetRows()
	if err != nil {
		return err
	}

	lastId := ""
	err = iterateChangeSetRows(d.changeLog, changeSetRows, func(before []*refactor.ChangeSet, csr *ChangeSetRow) error {
		if len(before) > 0 {
			return &ErrChangeSetOutOfOrder{
				ChangeSets:   before,
				ChangeSetRow: csr,
				BeforeRow:    true,
			}
		}

		lastId = csr.Id

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println(lastId)

	// changeSets := d.changeLog.Chan

	return nil
}
