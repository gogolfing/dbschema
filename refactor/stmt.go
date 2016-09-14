package refactor

type Stmt string

func (s Stmt) String() string {
	return string(s)
}

type StmtsFunc func(ctx Context) ([]Stmt, error)

func (f StmtsFunc) Validated(v Validator, ctx Context) ([]Stmt, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return f(ctx)
}

type StmtFunc func(ctx Context) (Stmt, error)

func StmtsFromFuncs(ctx Context, funcs ...StmtFunc) ([]Stmt, error) {
	stmts := []Stmt{}
	for _, f := range funcs {
		stmt, err := f(ctx)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}
