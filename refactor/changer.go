package refactor

type Changer interface {
	Stmts(ctx Context) (stmts []*Stmt, err error)
}
