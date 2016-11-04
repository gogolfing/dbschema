package refactor

import "fmt"

type Stmt struct {
	Raw    string
	Params []interface{}
}

func NewStmtFmt(format string, a ...interface{}) *Stmt {
	return NewStmt(fmt.Sprintf(format, a...))
}

func NewStmt(raw string, params ...interface{}) *Stmt {
	return &Stmt{
		Raw:    raw,
		Params: params,
	}
}

func (s *Stmt) AppendFmt(format string, a ...interface{}) *Stmt {
	return s.AppendRaw(fmt.Sprintf(format, a...))
}

func (s *Stmt) AppendRaw(raw string) *Stmt {
	return s.Append(NewStmt(raw))
}

func (s *Stmt) AppendParams(params ...interface{}) *Stmt {
	return s.Append(NewStmt("", params...))
}

func (s *Stmt) Append(others ...*Stmt) *Stmt {
	raw := s.Raw
	params := make([]interface{}, 0, len(s.Params))
	copy(params, s.Params)

	for _, other := range others {
		raw += other.Raw
		params = append(params, other.Params...)
	}

	return &Stmt{
		Raw:    raw,
		Params: params,
	}
}

func (s *Stmt) String() string {
	return fmt.Sprintf("%v %v", s.Raw, s.Params)
}

type StmtsFunc func(ctx Context) ([]*Stmt, error)

func (f StmtsFunc) Validated(v Validator, ctx Context) ([]*Stmt, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return f(ctx)
}

type StmtFunc func(ctx Context) (*Stmt, error)

func StmtsFromFuncs(ctx Context, funcs ...StmtFunc) ([]*Stmt, error) {
	stmts := []*Stmt{}
	for _, f := range funcs {
		stmt, err := f(ctx)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}
