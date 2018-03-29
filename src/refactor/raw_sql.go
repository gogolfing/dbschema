package refactor

import "github.com/gogolfing/dbschema/src/refactor/dto"

type RawSql struct {
	UpStmts   []*Stmt
	DownStmts []*Stmt

	dto *dto.RawSql
}

func (r *RawSql) Up(_ Context) ([]*Stmt, error) {
	//TODO: need to expand variables in the raw strings
	return r.UpStmts, nil
}

func (r *RawSql) Down(_ Context) ([]*Stmt, error) {
	//TODO: need to expand variables in the raw strings
	return r.DownStmts, nil
}

func (r *RawSql) DTO() (dto.Changer, error) {
	return r.dto, nil
}
