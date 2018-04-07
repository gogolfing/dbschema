package dbschema

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/gogolfing/dbschema/src/refactor"
)

const (
	DefaultChangeLogTableName     = "dbschema_changelog"
	DefaultChangeLogLockTableName = "dbschema_changelog_lock"
)

const DefaultLockId = "LOCK"

func (d *DBSchema) init() error {
	refCtx := d.initialRefactorContext()
	expanded, err := refactor.ExpandAll(
		refCtx,
		d.changeLog.TableName.DefaultExpander(DefaultChangeLogTableName),
		d.changeLog.LockTableName.DefaultExpander(DefaultChangeLogLockTableName),
	)
	if err != nil {
		return err
	}
	d.tableName, d.lockTableName = expanded[0], expanded[1]

	if err := d.initTables(); err != nil {
		return err
	}
	if err := d.initLock(); err != nil {
		return err
	}
	return nil
}

func (d *DBSchema) initTables() error {
	changeLogTable := createChangeLogTable(d.tableName)
	changeLogLockTable := createChangeLogLockTable(d.lockTableName)
	stmts := []*refactor.Stmt{
		changeLogTable,
		changeLogLockTable,
	}

	return d.executeNewTxWork(func(qe QueryExecer) error {
		return executeStmts(qe, stmts)
	})
}

type lock struct {
	id       string
	isLocked bool
	lockedAt time.Time
	lockedBy string
}

type errAlreadyLocked []*lock

func (e errAlreadyLocked) Error() string {
	lockResult := []string{}
	for _, lock := range []*lock(e) {
		lockResult = append(
			lockResult,
			fmt.Sprintf(
				"At %v by %v",
				lock.lockedAt.Format(DefaultTimeFormat),
				lock.lockedBy,
			),
		)
	}
	return fmt.Sprintf("dbschema: already locked...\n%v", strings.Join(lockResult, "\n"))
}

func (d *DBSchema) initLock() error {
	work := func(qe QueryExecer) error {
		locks, err := d.obtainLocks(qe)
		if err != nil {
			return err
		}

		if len(locks) > 0 {
			return errAlreadyLocked(locks)
		}

		return d.setLock(qe)
	}
	return d.executeNewTxWork(work)
}

func (d *DBSchema) obtainLocks(q Queryer) ([]*lock, error) {
	rows, err := q.Query(
		fmt.Sprintf(
			"%v WHERE %v = 1 AND %v = %v",
			d.selectFrom(
				d.lockTableName,
				ColumnLockId,
				ColumnIsLocked,
				ColumnLockedAt,
				ColumnLockedBy,
			),
			d.QuoteRef(ColumnIsLocked),
			d.QuoteRef(ColumnLockId),
			d.QuoteConst(DefaultLockId),
		),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locks := []*lock{}
	for rows.Next() {
		lock := &lock{}
		if err := scanLock(lock, rows); err != nil {
			return nil, err
		}
		locks = append(locks, lock)
	}
	return locks, nil
}

func scanLock(l *lock, s Scanner) error {
	isLockedInt := 0
	if err := s.Scan(&l.id, &isLockedInt, &l.lockedAt, &l.lockedBy); err != nil {
		return err
	}
	l.isLocked = isLockedInt == 1
	return nil
}

func (d *DBSchema) setLock(e Execer) error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	lockedBy := fmt.Sprintf("%v@%v", user.Username, hostname)

	_, err = e.Exec(
		fmt.Sprintf(
			"INSERT INTO %v VALUES (%v, %v, %v, %v)",
			d.QuoteRef(d.lockTableName),
			d.Placeholder(0),
			d.Placeholder(1),
			d.Placeholder(2),
			d.Placeholder(3),
		),
		DefaultLockId,
		"1",
		time.Now().UTC(),
		lockedBy,
	)
	return err
}
