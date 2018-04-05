package refactor

//RawSQL is a Changer that provides raw string sql statements with support for
//expansion of those raw input values.
type RawSQL struct {
	//UpStmts are the raw statements to return from Up after expansion.
	UpStmts []*Stmt

	//DownStmts are the raw statements to return from Down after expansion.
	DownStmts []*Stmt
}

//Up is the Changer implementation.
//It returns the expanded versions of r.UpStmts.
func (r *RawSQL) Up(ctx Context) ([]*Stmt, error) {
	result := make([]*Stmt, 0, len(r.UpStmts))
	for i, upStmt := range r.UpStmts {
		expanded, err := ctx.Expand(upStmt.Raw)
		if err != nil {
			return nil, err
		}
		result[i] = &Stmt{Raw: expanded}
	}
	return result, nil
}

//Down is the Changer implementation.
//It returns the expanded versions of r.DownStmts.
func (r *RawSQL) Down(ctx Context) ([]*Stmt, error) {
	result := make([]*Stmt, 0, len(r.DownStmts))
	for i, downStmt := range r.DownStmts {
		expanded, err := ctx.Expand(downStmt.Raw)
		if err != nil {
			return nil, err
		}
		result[i] = &Stmt{Raw: expanded}
	}
	return result, nil
}
