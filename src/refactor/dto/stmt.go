package dto

import (
	"encoding/xml"

	"github.com/gogolfing/dbschema/src/refactor"
)

type Stmt struct {
	XMLName xml.Name `xml:"Stmt"`

	Raw string `xml:",innerxml"`
}

func (s *Stmt) RefactorType() *refactor.Stmt {
	return &refactor.Stmt{
		Raw: s.Raw,
	}
}
