package dbschema

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/gogolfing/dbschema/logger"
	"github.com/gogolfing/dbschema/refactor"
)

type ChangeSetRow struct {
	Id            string
	Name          *string
	Author        *string
	ExecutedAt    time.Time
	UpdatedAt     time.Time
	OrderExecuted int
	Sha256Sum     string
	Tags          []string
	Version       string
}

func (csr *ChangeSetRow) String() string {
	nameAuthor := ""
	if csr.Name != nil {
		nameAuthor += *csr.Name
	}
	if csr.Author != nil {
		if nameAuthor != "" && *csr.Author != "" {
			nameAuthor += " by "
		}
		nameAuthor += *csr.Author
	}
	if nameAuthor != "" {
		nameAuthor += "\n"
	}

	at := fmt.Sprintf("Executed at %v", csr.ExecutedAt.Format(DefaultTimeFormat))
	if csr.UpdatedAt != csr.ExecutedAt {
		at += fmt.Sprintf(" and updated at %v", csr.UpdatedAt.Format(DefaultTimeFormat))
	}

	return fmt.Sprintf(
		"ChangeSet - %v\n%v%v",
		csr.Id,
		nameAuthor,
		at,
	)
}

func (csr *ChangeSetRow) StringVerbose() string {
	tags := make([]string, 0, len(csr.Tags))
	for _, tag := range csr.Tags {
		tags = append(tags, fmt.Sprintf("%q", tag))
	}

	return fmt.Sprintf(
		`SHA-256 Sum:      %v
Tags:             %v
DBSchema Version: %v`,
		csr.Sha256Sum,
		strings.Join(tags, ", "),
		csr.Version,
	)
}

func (d *DBSchema) listOrderedChangeSetRows() ([]*ChangeSetRow, error) {
	result := []*ChangeSetRow{}

	work := func(qe QueryExecer) error {
		rows, err := qe.Query(
			refactor.NewStmtFmt(
				"%v\nORDER BY %v ASC",
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
				d.QuoteRef(ColumnOrderExecuted),
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

func printChangeSetRow(logger logger.Logger, csr *ChangeSetRow) {
	fmt.Fprintln(logger.Info(), csr)
	fmt.Fprintln(logger.Verbose(), csr.StringVerbose())
}
