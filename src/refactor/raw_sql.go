package refactor

type RawSql struct {
	UpStmts   []*Stmt
	DownStmts []*Stmt
}

func (r *RawSql) Up(_ Context) ([]*Stmt, error) {
	//TODO: need to expand variables in the raw strings
	return r.UpStmts, nil
}

func (r *RawSql) Down(_ Context) ([]*Stmt, error) {
	//TODO: need to expand variables in the raw strings
	return r.DownStmts, nil
}
