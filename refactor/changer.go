package refactor

type Changer interface {
	Up(ctx Context) (stmts []Stmt, err error)
	Down(ctx Context) (stmts []Stmt, err error)
}
