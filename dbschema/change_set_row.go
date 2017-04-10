package dbschema

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
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

const changeSetRowLineFormat = "%-16s%v\n"

func (csr *ChangeSetRow) String() string {
	buffer := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buffer, "ChangeSet - %s\n", csr.Id)
	if csr.Name != nil {
		fmt.Fprintf(buffer, changeSetRowLineFormat, "Name:", *csr.Name)
	}
	if csr.Author != nil {
		fmt.Fprintf(buffer, changeSetRowLineFormat, "Author:", *csr.Author)
	}
	fmt.Fprintf(buffer, changeSetRowLineFormat, "Executed At:", csr.ExecutedAt.Format(DefaultTimeFormat))

	return strings.TrimRight(buffer.String(), "\n")
}

func (csr *ChangeSetRow) StringVerbose() string {
	buffer := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buffer, changeSetRowLineFormat, "SHA-256 Sum:", csr.Sha256Sum)
	tags := make([]string, 0, len(csr.Tags))
	for _, tag := range csr.Tags {
		tags = append(tags, fmt.Sprintf("%q", tag))
	}
	tagOutput := "<none>"
	if len(tags) > 0 {
		tagOutput = strings.Join(tags, ", ")
	}
	fmt.Fprintf(buffer, changeSetRowLineFormat, "Tags:", tagOutput)

	return strings.TrimRight(buffer.String(), "\n")
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
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, &ParsingTagsError{
			Tags: tags,
			Err:  err,
		}
	}
	return row, nil
}

func printChangeSetRow(logger logger.Logger, csr *ChangeSetRow) {
	fmt.Fprintln(logger.Info(), csr)
	fmt.Fprintln(logger.Verbose(), csr.StringVerbose())
	fmt.Fprintln(logger.Info())
}
