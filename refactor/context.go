package refactor

import "github.com/gogolfing/dbschema/dialect"

type Context interface {
	dialect.Dialect
	Expand(expr string) (value string, err error)
}

func ExpandAll(ctx Context, expanders ...Expander) ([]string, error) {
	result := []string{}
	for _, expander := range expanders {
		value, err := expander.Expand(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

type Expander interface {
	Expand(Context) (string, error)
}

type ExpanderFunc func(Context) (string, error)

func (ef ExpanderFunc) Expand(ctx Context) (string, error) {
	return ef(ctx)
}

func StringExpander(s *string, def string) Expander {
	return ExpanderFunc(func(ctx Context) (string, error) {
		if s == nil {
			return def, nil
		}
		return ctx.Expand(*s)
	})
}
