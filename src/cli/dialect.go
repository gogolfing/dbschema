package cli

import (
	"fmt"

	"github.com/gogolfing/dbschema/src/dialect"
	"github.com/gogolfing/dbschema/src/dialect/mysql"
	"github.com/gogolfing/dbschema/src/dialect/postgresql"
)

type UnsupportedDBMSError string

func (e UnsupportedDBMSError) Error() string {
	return fmt.Sprintf("cli: unsupported dbms %s", string(e))
}

func NewDialect(dbms string) (dialect.Dialect, error) {
	var result dialect.Dialect

	switch dbms {
	case postgresql.DBMS:
		result = postgresql.Dialect()

	case mysql.DBMS:
		result = mysql.Dialect()
	}

	if result == nil {
		return nil, UnsupportedDBMSError(dbms)
	}

	return result, nil
}
