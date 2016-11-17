package dbschema

import (
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/refactor"
)

func collectingAppliedChangeSets(w io.Writer) {
	fmt.Fprintln(w, "Collecting applied ChangeSets...\n")
}

func iterateChangeSetRows(changeLog *refactor.ChangeLog, csrows []*ChangeSetRow, callback func(before []*refactor.ChangeSet, csr *ChangeSetRow) error) error {
	prevId := ""
	for _, csr := range csrows {
		changeSetsBefore := changeLog.ChangeSetsSubSlice(prevId, csr.Id)

		err := callback(changeSetsBefore, csr)
		if err != nil {
			return err
		}

		prevId = csr.Id
	}

	return nil
}
