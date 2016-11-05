package dbschema

import (
	"time"

	"github.com/gogolfing/dbschema/refactor"
)

type ChangeSetRow struct {
	Id            string
	Name          string
	Author        string
	ExecutedAt    time.Time
	UpdatedAt     time.Time
	OrderExecuted int
	Sha256Sum     string
	Tags          []string
	Version       string
}

func (d *DBSchema) collectChangeSetRows() ([]*ChangeSetRow, error) {
	result := []*ChangeSetRow{}

	work := func(qe QueryExecer) error {
		rows, err := qe.Query(
			refactor.NewStmtFmt(
				d.selectFrom(
					d.changeLogTableName,
					ColumnChangeSetId,
				),
			),
		)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			csr := &ChangeSetRow{}
			if err := scanChangeSetRow(csr, rows); err != nil {
				return err
			}
			result = append(result, csr)
		}

		return nil
	}

	if err := d.executeTxWork(work); err != nil {
		return nil, err
	}
	return result, nil
}

func scanChangeSetRow(csr *ChangeSetRow, s Scanner) error {
	return s.Scan(
		&csr.Id,
	)
}
