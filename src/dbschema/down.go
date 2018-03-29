package dbschema

import (
	"github.com/gogolfing/dbschema/src/logger"
	"github.com/gogolfing/dbschema/src/refactor"
)

func (d *DBSchema) Down(logger logger.Logger, count int) error {
	collectingAppliedChangeSets(logger.Verbose())

	err := d.ensureAppliedChangeSets(true, nil)
	if err != nil {
		return err
	}

	appliedSets, err := d.listOrderedAppliedChangeSets()
	if err != nil {
		return err
	}

	changeSets := d.changeLog.ChangeSets[:len(appliedSets)]

	work := func(qe QueryExecer) error {
		for i := len(changeSets) - 1; i >= 0; i-- {
			changeSet := changeSets[i]

			stmts, err := refactor.CollectChangersDown(d, changeSet.Changers...)
			if err != nil {
				return err
			}

			//TODO: send a logger here.
			if err := executeStmts(qe, stmts); err != nil {
				return err
			}

			//if not dry run.
			if err := d.deleteFromChangeLogTable(qe, changeSet); err != nil {
				return err
			}
		}
		return nil
	}

	return d.executeNewTxWork(work)
}
