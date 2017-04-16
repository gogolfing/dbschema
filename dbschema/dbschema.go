package dbschema

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/refactor"
)

const DefaultTimeFormat = time.RFC1123Z

type DBSchema struct {
	db DB

	dialect.Dialect

	changeLog *refactor.ChangeLog

	tableName     string
	lockTableName string
}

func Open(dialect dialect.Dialect, conn string, changeLog *refactor.ChangeLog) (*DBSchema, error) {
	sqlDB, err := openSqlDB(dialect, conn)
	if err != nil {
		return nil, err
	}
	return openDB(sqlDB, dialect, changeLog)
}

func openDB(db DB, dialect dialect.Dialect, changeLog *refactor.ChangeLog) (*DBSchema, error) {
	d := &DBSchema{
		db:        db,
		Dialect:   dialect,
		changeLog: changeLog,
	}
	if err := d.open(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DBSchema) open() error {
	if err := d.db.Ping(); err != nil {
		return err
	}
	if err := d.init(); err != nil {
		return err
	}
	return nil
}

func (d *DBSchema) Close() error {
	if err := d.finalize(); err != nil {
		return err
	}
	return d.db.Close()
}

func (d *DBSchema) Expand(expr string) (value string, err error) {
	return dialect.Expand(expr, d.changeLog.Variables, d.Dialect)
}

func (d *DBSchema) executeNewTxWork(work func(QueryExecer) error) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	err = work(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d *DBSchema) selectFrom(table string, columns ...string) string {
	columnRefs := make([]string, 0, len(columns))
	for _, column := range columns {
		columnRefs = append(columnRefs, d.QuoteRef(column))
	}
	return fmt.Sprintf(
		"SELECT %v FROM %v",
		strings.Join(columnRefs, ", "),
		d.QuoteRef(table),
	)
}

func (d *DBSchema) listOrderedAppliedChangeSets() ([]*AppliedChangeSet, error) {
	result := []*AppliedChangeSet{}

	work := func(qe QueryExecer) error {
		rows, err := qe.Query(
			fmt.Sprintf(
				"%v ORDER BY %v ASC",
				d.selectFrom(
					d.tableName,
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
			acs := &AppliedChangeSet{}
			if err := scanChangeSetRow(acs, rows); err != nil {
				return err
			}
			result = append(result, acs)
		}

		return nil
	}

	if err := d.executeNewTxWork(work); err != nil {
		return nil, err
	}
	return result, nil
}

func scanChangeSetRow(acs *AppliedChangeSet, s Scanner) error {
	tagsRaw := sql.NullString{}
	err := s.Scan(
		&acs.Id,
		&acs.Name,
		&acs.Author,
		&acs.ExecutedAt,
		&acs.UpdatedAt,
		&acs.OrderExecuted,
		&acs.Sha256Sum,
		&tagsRaw,
		&acs.Version,
	)
	if err != nil {
		return err
	}

	if tagsRaw.Valid {
		tags, err := parseTags(tagsRaw.String)
		if err != nil {
			return err
		}
		acs.Tags = tags
	}

	return nil
}

func (d *DBSchema) ensureAppliedChangeSets(fix bool, visitor func(acs *AppliedChangeSet)) error {
	appliedSets, err := d.listOrderedAppliedChangeSets()
	if err != nil {
		return err
	}

	for i, acs := range appliedSets {
		if i >= len(d.changeLog.ChangeSets) {
			return fmt.Errorf("applied change set not in change log")
		}
		changeSet := d.changeLog.ChangeSets[i]
		if changeSet.Id != acs.Id {
			return fmt.Errorf("applied change set does not equal change set in change log at index %d", i)
		}

		if visitor != nil {
			visitor(acs)
		}
	}

	return nil
}

func (d *DBSchema) insertIntoChangeLogTable(e Execer, changeSet *refactor.ChangeSet, order int) error {
	now := time.Now().UTC()
	hash, err := changeSet.Sha256Sum()
	if err != nil {
		return err
	}
	tags, tagsOk := formatTags(changeSet.Tags)
	query := fmt.Sprintf("INSERT INTO %s VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)",
		d.tableName,
		d.Placeholder(0),
		d.Placeholder(1),
		d.Placeholder(2),
		d.Placeholder(3),
		d.Placeholder(4),
		d.Placeholder(5),
		d.Placeholder(6),
		d.Placeholder(7),
		d.Placeholder(8),
	)
	_, err = e.Exec(
		query,
		changeSet.Id,
		sql.NullString{String: changeSet.Name.String, Valid: changeSet.Name.Valid},
		sql.NullString{String: changeSet.Author.String, Valid: changeSet.Author.Valid},
		now,
		now,
		order,
		hex.EncodeToString(hash),
		sql.NullString{String: tags, Valid: tagsOk},
		Version,
	)
	return err
}

func (d *DBSchema) deleteFromChangeLogTable(e Execer, changeSet *refactor.ChangeSet) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		d.tableName,
		ColumnChangeSetId,
		d.Placeholder(0),
	)
	_, err := e.Exec(query, changeSet.Id)
	return err
}

//TODO: add a logger here.
func executeStmts(qe QueryExecer, stmts []*refactor.Stmt) error {
	for _, stmt := range stmts {
		_, err := qe.Exec(stmt.Raw, stmt.Args...)
		if err != nil {
			return err
		}
	}
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

func formatTags(tags []string) (string, bool) {
	if len(tags) == 0 {
		return "", false
	}
	buffer := bytes.NewBuffer([]byte{})
	w := csv.NewWriter(buffer)
	w.Write(tags)
	w.Flush()
	return buffer.String(), true
}
