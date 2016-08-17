package refactor

import (
	"encoding/xml"
	"fmt"
)

type AddTable struct {
	XMLName xml.Name `xml:"AddTable"`

	Name string `xml:"name,attr"`

	Columns []*Column `xml:"Columns"`
}

func (a *AddTable) Validate() error {
	if a.Name == "" {
		return ValidationError("AddTable.Name cannot be empty")
	}
	return nil
}

func (a *AddTable) Up(ctx *Context) (stmts []string, err error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return StmtsFromFuncs(ctx, a.upTable)
}

func (a *AddTable) upTable(ctx *Context) (string, error) {
	columns, err := a.upColumns(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v %v (%v)", ctx.CreateTable, a.Name, columns), nil
}

func (a *AddTable) upColumns(ctx *Context) (string, error) {
	return "", nil
}
