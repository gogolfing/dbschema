package dbschema

import "github.com/gogolfing/dbschema/refactor"

func (d *DBSchema) finalize() error {
	return d.finalizeLock()
}

func (d *DBSchema) finalizeLock() error {
	work := func(qe QueryExecer) error {
		return d.unsetLock(qe)
	}
	return d.executeTxWork(work)
}

func (d *DBSchema) unsetLock(e Execer) error {
	_, err := e.Exec(
		refactor.NewStmtFmt(
			"DELETE FROM %v WHERE %v = %v AND %v = %v",
			d.QuoteRef(d.changeLogLockTableName),
			d.QuoteRef(ColumnLockId), d.Placeholder(0),
			d.QuoteRef(ColumnIsLocked), d.Placeholder(1),
		).AppendParams(
			DefaultLockId,
			1,
		),
	)
	return err
}
