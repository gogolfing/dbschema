package refactor

import "github.com/gogolfing/dbschema/dialect"

type Context interface {
	dialect.Dialect
	Expand(expr string) (value string, err error)
}

func ExpandAll(ctx Context, expressions ...string) ([]string, error) {
	result := []string{}
	for _, expr := range expressions {
		if value, err := ctx.Expand(expr); err != nil {
			return nil, err
		} else {
			result = append(result, value)
		}
	}
	return result, nil
}
