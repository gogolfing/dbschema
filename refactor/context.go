package refactor

import "github.com/gogolfing/dbschema/dialect"

type Context interface {
	dialect.Dialect
	Expand(expr string) (value string, err error)
}
