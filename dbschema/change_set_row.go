package dbschema

import (
	"encoding/csv"
	"strings"
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
					ColumnChangeSetName,
					ColumnChangeSetAuthor,
					ColumnExecutedAt,
					ColumnUpdatedAt,
					ColumnOrderExecuted,
					ColumnSha256Sum,
					ColumnTags,
					ColumnVersion,
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
	tagsRaw := ""
	err := s.Scan(
		&csr.Id,
		&csr.Name,
		&csr.Author,
		&csr.ExecutedAt,
		&csr.UpdatedAt,
		&csr.OrderExecuted,
		&csr.Sha256Sum,
		&tagsRaw,
		&csr.Version,
	)
	if err != nil {
		return err
	}
	tags, err := parseTags(tagsRaw)
	if err != nil {
		return err
	}
	csr.Tags = tags
	return nil
}

func parseTags(tags string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(tags))
	row, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return row, nil
}
