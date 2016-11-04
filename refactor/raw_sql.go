package refactor

import "encoding/xml"

const (
	errMustBeNonEmpty = ErrInvalid("RawSql > Stmt(s) must be non empty")
)

type RawSql struct {
	XMLName xml.Name `xml:"RawSql"`

	StmtSlice []*Stmt `xml:"Stmt"`
}

func (r *RawSql) Validate() error {
	if len(r.StmtSlice) == 0 {
		return errMustBeNonEmpty
	}
	return nil
}

func (r *RawSql) Stmts(ctx Context) ([]*Stmt, error) {
	return StmtsFunc(r.stmts).Validated(r, ctx)
}

func (r *RawSql) stmts(ctx Context) ([]*Stmt, error) {
	return r.StmtSlice, nil
}
