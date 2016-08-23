package refactor

import (
	"github.com/gogolfing/dbschema/dialect"
	_vars "github.com/gogolfing/dbschema/vars"
)

const (
	DefaultTestContextName = "dbscema/refactor/context_test.go/test"
	DefaultTestContextDBMS = "dbscema/refactor/context_test.go/dbms"
)

func newDefaultTestContext() Context {
	return newTestContext(nil, "", nil, []string{}...)
}

func newTestContext(d dialect.Dialect, dbms string, names []string, vars ...string) Context {
	if dbms == "" {
		dbms = dbms
	}
	if len(names) == 0 {
		names = []string{DefaultTestContextName}
	}
	if len(vars)%2 == 1 {
		vars = vars[:len(vars)-1]
	}
	v := &_vars.Variables{}
	for i := 0; i < len(vars); i += 2 {
		v.Put(&_vars.Variable{Name: vars[i], Value: vars[i+1]})
	}
	return &context{dialect: d, names: names, dbms: dbms, vars: v}
}
