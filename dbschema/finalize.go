package dbschema

import "fmt"

func (d *DBSchema) finalize() error {
	return d.finalizeLock()
}

func (d *DBSchema) finalizeLock() error {
	work := func(qe QueryExecer) error {
		return d.unsetLock(qe)
	}
	return d.executeNewTxWork(work)
}

func (d *DBSchema) unsetLock(e Execer) error {
	_, err := e.Exec(
		fmt.Sprintf(
			"DELETE FROM %v WHERE %v = %v AND %v = %v",
			d.QuoteRef(d.lockTableName),
			d.QuoteRef(ColumnLockId), d.Placeholder(0),
			d.QuoteRef(ColumnIsLocked), d.Placeholder(1),
		),
		DefaultLockId,
		1,
	)
	return err
}
