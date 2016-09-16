package refactor

import (
	"encoding/xml"
	"fmt"

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
		return ErrInvalid("CreateTable.name cannot be empty")
	}
	if err := strval.ValidateBool(c.IfNotExists); err != nil {
		return fmt.Errorf("CreateTable.ifNotExists %v", err)
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
	definition, err := c.definition(ctx)
	if err != nil {
		return nil, err
	}
	constraints, err := c.constraints(ctx)
	if err != nil {
		return nil, err
	}
	return append([]Stmt{definition}, constraints...), nil
}

func (c *CreateTable) definition(ctx Context) (Stmt, error) {
	expanded, err := ExpandAll(ctx, c.Name)
	if err != nil {
		return "", err
	}
	name := expanded[0]

	result := dialect.CreateTable
	if strval.Bool(c.IfNotExists, false) {
		result += " " + dialect.IfNotExists
	}
	result += " " + ctx.QuoteRef(name) + " (\n"

	colDefs, err := c.columnDefinitions(ctx)
	if err != nil {
		return "", err
	}
	result += colDefs + "\n)"

	return Stmt(result), nil
}

func (c *CreateTable) columnDefinitions(ctx Context) (string, error) {
	result := ""
	for i, col := range c.Columns {
		colDef, err := col.Definition(ctx)
		if err != nil {
			return "", err
		}
		result += "\t" + colDef
		if i != len(c.Columns)-1 {
			result += ",\n"
		}
	}
	return result, nil
}

func (c *CreateTable) constraints(ctx Context) ([]Stmt, error) {
	return nil, nil
}

func (c *CreateTable) Down(ctx Context) ([]Stmt, error) {
	return nil, nil
}
