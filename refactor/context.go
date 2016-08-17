package refactor

import (
	"fmt"
	"strings"

	"github.com/gogolfing/dbschema/dialect"
)

const (
	VariablePrefix = "{"
	VariableSuffix = "}"
)

const DialectVariablePrefix = "Dialect."

type ErrVariableDoesNotExist string

func (e ErrVariableDoesNotExist) Error() string {
	return fmt.Sprintf("refactor: variable does not exist %q", string(e))
}

type ErrInvalidVariableReference string

func (e ErrInvalidVariableReference) Error() string {
	return fmt.Sprintf("refactor: invalid variable reference %q", string(e))
}

type Context struct {
	*dialect.Dialect

	Variables map[string]string
}

func NewContext(d *dialect.Dialect) *Context {
	if d == nil {
		d = dialect.NewSqlDialect()
	}
	return &Context{
		Dialect: d,

		Variables: map[string]string{},
	}
}

func (c *Context) GetVariable(name string) (string, error) {
	if !strings.HasPrefix(name, VariablePrefix) || !strings.HasSuffix(name, VariableSuffix) {
		return "", ErrInvalidVariableReference(name)
	}
	name = name[1 : len(name)-1]
	origName := name

	value, ok := c.Variables[name]
	if !ok {
		if !strings.HasPrefix(name, DialectVariablePrefix) {
			return "", ErrVariableDoesNotExist(origName)
		}
		name = strings.TrimPrefix(name, DialectVariablePrefix)
		value, err := c.Dialect.ValueOfVariableField(name)
		if err != nil {
			return "", ErrVariableDoesNotExist(origName)
		}
		return value, nil
	}
	return value, nil
}
