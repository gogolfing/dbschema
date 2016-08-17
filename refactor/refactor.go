package refactor

import "fmt"

const (
	True = "true"
)

func BoolDefault(attr string, def bool) bool {
	if attr == "" {
		return def
	}
	return attr == True
}

type Validator interface {
	Validate() error
}

type ValidationError string

func (e ValidationError) Error() string {
	return fmt.Sprintf("refactor: validation error: %v", string(e))
}

func StmtsFromFuncs(ctx *Context, stmtFuncs ...func(ctx *Context) (stmt string, err error)) (stmts []string, err error) {
	for _, stmtFunc := range stmtFuncs {
		stmt, err := stmtFunc(ctx)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}
