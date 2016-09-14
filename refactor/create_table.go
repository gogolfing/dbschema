package refactor

import (
	"encoding/xml"

	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/refactor/strval"
)

type CreateTable struct {
	XMLName xml.Name `xml:"CreateTable"`

	Name string `xml:"name,attr"`

	IfNotExists *string `xml:"ifNotExists,attr"`

	Columns []*Column `xml:"Columns"`
}

func (c *CreateTable) Validate() error {
	if c.Name == "" {
		return ErrInvalid("CreateTable.Name cannot be empty")
	}
	if err := strval.ValidateBool(c.IfNotExists); err != nil {
		return err
	}
	for _, col := range c.Columns {
		if err := col.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CreateTable) Up(ctx Context) ([]Stmt, error) {
	return StmtsFunc(c.up).Validated(c, ctx)
}

func (c *CreateTable) up(ctx Context) ([]Stmt, error) {
	result := dialect.CreateTable
	return []Stmt{Stmt(result)}, nil
}

func (c *CreateTable) Down(ctx Context) ([]Stmt, error) {
	return nil, nil
}
