package refactor

type Changer interface {
	Up(ctx Context) (stmts []string, err error)
	Down(ctx Context) (stmts []string, err error)

	Validator
}
