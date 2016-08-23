package refactor

import (
	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/vars"
)

type Context interface {
	Dialect() dialect.Dialect

	DBMS() string

	Expand(expr string) (value string, err error)

	Names() []string
}

type context struct {
	dialect dialect.Dialect

	dbms string

	names []string

	vars *vars.Variables
}

func (c *context) Dialect() dialect.Dialect {
	return c.dialect
}

func (c *context) DBMS() string {
	return c.dbms
}

func (c *context) Expand(expr string) (value string, err error) {
	return c.vars.Dereference(expr)
}

func (c *context) Names() []string {
	return c.names
}
