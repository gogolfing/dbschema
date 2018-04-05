package dto

import (
	"encoding/xml"

	"github.com/gogolfing/dbschema/src/refactor"
)

type RawSQL struct {
	XMLName xml.Name `xml:"RawSQL"`

	Up *RawSQLStmts `xml:"Up"`

	Down *RawSQLStmts `xml:"Down"`
}

func TransformRawSQL(c *refactor.RawSQL) *RawSQL {
	return &RawSQL{
		Up:   TransformSQLStmts(c.UpStmts),
		Down: TransformSQLStmts(c.DownStmts),
	}
}

type RawSQLStmts struct {
	Stmts []*Stmt `xml:"Stmt"`
}

func TransformSQLStmts(stmts []*refactor.Stmt) *RawSQLStmts {
	result := &RawSQLStmts{}
	for _, stmt := range stmts {
		result.Stmts = append(
			result.Stmts,
			&Stmt{
				Raw: stmt.Raw,
			},
		)
	}
	return result
}
