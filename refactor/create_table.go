package refactor

import (
	"encoding/xml"
	"strings"

	"github.com/gogolfing/dbschema/dialect"
)

type CreateTable struct {
	XMLName xml.Name `xml:"CreateTable"`

	Name *StringAttr `xml:"name,attr"`

	IfNotExists *BoolAttr `xml:"ifNotExists,attr"`

	Columns []*Column `xml:"Column"`
}

func (c *CreateTable) Validate() error {
	return ValidateAll(
		c.Name.NotEmptyValidator("CreateTable.name"),
		c.IfNotExists.Validator("CreateTable.ifNotExists"),
		ColumnsValidator(c.Columns),
	)
}

func (c *CreateTable) Stmts(ctx Context) ([]*Stmt, error) {
	return StmtsFunc(c.stmts).Validated(c, ctx)
}

func (c *CreateTable) stmts(ctx Context) ([]*Stmt, error) {
	definition, err := c.definition(ctx)
	if err != nil {
		return nil, err
	}
	constraints, err := c.constraints(ctx)
	if err != nil {
		return nil, err
	}
	return append([]*Stmt{definition}, constraints...), nil
}

func (c *CreateTable) definition(ctx Context) (*Stmt, error) {
	expanded, err := ExpandAll(
		ctx,
		c.Name,
		c.IfNotExists.Expander(false),
	)
	if err != nil {
		return nil, err
	}
	name, ifNotExists := expanded[0], BoolString(expanded[1])

	result := dialect.CreateTable

	if ifNotExists {
		result += " " + dialect.IfNotExists
	}

	result += " " + ctx.QuoteRef(name) + " ("

	colDefs, err := c.columnDefinitions(ctx)
	if err != nil {
		return nil, err
	}
	if colDefs != "" {
		result += "\n" + colDefs + "\n"
	}
	result += ")"

	return NewStmt(result), nil
}

func (c *CreateTable) columnDefinitions(ctx Context) (string, error) {
	colDefs := []string{}
	for _, col := range c.Columns {
		colDef, err := col.Definition(ctx)
		if err != nil {
			return "", err
		}
		colDefs = append(colDefs, "\t"+colDef)
	}
	return strings.Join(colDefs, ",\n"), nil
}

func (c *CreateTable) constraints(ctx Context) ([]*Stmt, error) {
	result := []*Stmt{}

	pks, err := c.primaryKeyConstraints(ctx)
	if err != nil {
		return nil, err
	}
	result = append(result, pks...)

	return result, nil
}

func (c *CreateTable) primaryKeyConstraints(ctx Context) ([]*Stmt, error) {
	name := ""
	ics := []*IndexColumn{}

	for _, col := range c.Columns {
		tempName, ic, err := col.primaryKeyIndexColumn(ctx)
		if err != nil {
			return nil, err
		}
		if ic != nil {
			ics = append(ics, ic)
		}
		if name != "" && tempName != "" {
			name = tempName
		}
	}

	if len(ics) == 0 {
		return nil, nil
	}

	pk := &AddPrimaryKey{
		Table:        c.Name,
		Name:         NewStringAttr(name),
		IndexColumns: ics,
	}

	return pk.Stmts(ctx)
}
