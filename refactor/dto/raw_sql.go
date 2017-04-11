package dto

import "encoding/xml"

type RawSql struct {
	XMLName xml.Name `xml:"RawSql"`

	Up struct {
		XMLName xml.Name `xml:"Up"`

		Stmts []*Stmt `xml:"Stmt"`
	}

	Down struct {
		XMLName xml.Name `xml:"Down"`

		Stmts []*Stmt `xml:"Stmt"`
	}
}
