package dbschema

import (
	"fmt"
	"log"

	"github.com/gogolfing/dbschema/src/logger"
	"github.com/gogolfing/dbschema/src/refactor"
)

func (d *DBSchema) Up(logger logger.Logger, count int) error {
	collectingAppliedChangeSets(logger.Verbose())

	err := d.ensureAppliedChangeSets(true, nil)
	if err != nil {
		return err
	}

	appliedSets, err := d.listOrderedAppliedChangeSets()
	if err != nil {
		return err
	}
	lastOrderExecuted := 0
	if len(appliedSets) > 0 {
		lastOrderExecuted = appliedSets[len(appliedSets)-1].OrderExecuted
	}

	log.Println("d.Up()", d.changeLog.ChangeSets)

	refCtx := d.initialRefactorContext()
	toApply := d.changeLog.ChangeSets[len(appliedSets):]

	work := func(qe QueryExecer) error {
		for i, changeSet := range toApply {
			stmts, err := refactor.CollectChangersUp(refCtx, changeSet.Changers...)
			log.Println("stmts, err", stmts, err)
			if err != nil {
				return err
			}
			//TODO: send a logger here.

			fmt.Fprintln(logger.Verbose(), stmts)

			if err := executeStmts(qe, stmts); err != nil {
				return err
			}
			//if not dry run.
			if err := d.insertIntoChangeLogTable(qe, changeSet, lastOrderExecuted+i+1); err != nil {
				return err
			}
		}
		return nil
	}

	return d.executeNewTxWork(work)
}
