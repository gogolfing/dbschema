package refactor

import "github.com/gogolfing/dbschema/dialect"

type Context interface {
	ExpandErr(in string) (string, error)
}

func NewContext(d dialect.Dialect) *Context {
	return &Context{
		Dialect: d,
	}
}

func (c *Context) Expand(in string) string {

}

// func (c *Context) Expand(v string) string {
// 	value, err := c.GetVariableValue(v)
// 	if err != nil {
// 		return v
// 	}
// 	return value
// }

// func (c *Context) GetVariableValue(name string) (string, error) {
// 	if !strings.HasPrefix(name, VariablePrefix) || !strings.HasSuffix(name, VariableSuffix) {
// 		return "", ErrInvalidVariableReference(name)
// 	}
// 	name = name[1 : len(name)-1]
// 	origName := name

// 	value, ok := c.Variables[name]
// 	if !ok {
// 		if !strings.HasPrefix(name, DialectVariablePrefix) {
// 			return "", ErrVariableDoesNotExist(origName)
// 		}
// 		name = strings.TrimPrefix(name, DialectVariablePrefix)
// 		value, err := c.Dialect.ValueOfVariableField(name)
// 		if err != nil {
// 			return "", ErrVariableDoesNotExist(origName)
// 		}
// 		return value, nil
// 	}
// 	return value, nil
// }
