package refactor

import "encoding/xml"

const (
	errUpMustBeNonEmpty   = ErrInvalid("RawSql > Up > Stmt(s) must be non empty")
	errDownMustBeNonEmpty = ErrInvalid("RawSql > Down > Stmt(s) must be non empty")
)

type RawSql struct {
	XMLName xml.Name `xml:"RawSql"`

	UpStmts   []Stmt `xml:"Up>Stmt"`
	DownStmts []Stmt `xml:"Down>Stmt"`
}

func (r *RawSql) Validate(ctx Context) error {
	if len(r.UpStmts) == 0 {
		return errUpMustBeNonEmpty
	}
	if len(r.DownStmts) == 0 {
		return errDownMustBeNonEmpty
	}
	return nil
}

func (r *RawSql) Up(ctx Context) ([]Stmt, error) {
	return StmtsFunc(r.up).Validated(r, ctx)
}

func (r *RawSql) up(ctx Context) ([]Stmt, error) {
	return r.UpStmts, nil
}

func (r *RawSql) Down(ctx Context) ([]Stmt, error) {
	return StmtsFunc(r.down).Validated(r, ctx)
}

func (r *RawSql) down(ctx Context) ([]Stmt, error) {
	return r.DownStmts, nil
}
