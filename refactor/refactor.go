package refactor

import "fmt"

const (
	True  = "true"
	False = "false"
)

func StringDefault(value *string, def string) string {
	if value == nil {
		return def
	}
	return *value
}

func StringDefaultBool(value *string, def bool) bool {
	if value == nil {
		return def
	}
	return *value == True
}

func ValidateStringBool(value *string) error {
	if value == nil {
		return nil
	}
	if *value != True || *value != False {
		return fmt.Errorf("must be %q or %q", True, False)
	}
	return nil
}

type ErrInvalid string

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("refactor: invalid: %v", string(e))
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
