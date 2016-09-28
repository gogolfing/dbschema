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
	expanded, err := ExpandAll(
		ctx,
		c.Name,
		c.IfNotExists.Expander(false),
	)
	if err != nil {
		return "", err
	}
	name, ifNotExists := expanded[0], BoolString(expanded[1])

	result := dialect.CreateTable

	if ifNotExists {
		result += " " + dialect.IfNotExists
	}

	result += " " + ctx.QuoteRef(name) + " ("

	colDefs, err := c.columnDefinitions(ctx)
	if err != nil {
		return "", err
	}
	if colDefs != "" {
		result += "\n" + colDefs + "\n"
	}
	result += ")"

	return Stmt(result), nil
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

func (c *CreateTable) constraints(ctx Context) ([]Stmt, error) {
	result := []Stmt{}

	pk, err := c.primaryKeyConstraint(ctx)
	if err != nil {
		return nil, err
	}
	if pk != "" {
		result = append(result, pk)
	}

	return result, nil
}

func (c *CreateTable) primaryKeyConstraint(ctx Context) (Stmt, error) {
	name := ""
	ics := []*IndexColumn{}

	for _, col := range c.Columns {
		tempName, ic, err := col.primaryKeyIndexColumn(ctx)
		if err != nil {
			return "", err
		}
		if ic != nil {
			ics = append(ics, ic)
		}
		if name != "" && tempName != "" {
			name = tempName
		}
	}

	if len(ics) == 0 {
		return "", nil
	}

	pk := &AddPrimaryKey{
		Table:        c.Name,
		Name:         NewStringAttr(name),
		IndexColumns: ics,
	}

	stmts, err := pk.Up(ctx)
	if err != nil {
		return "", err
	}

	return stmts[0], nil
}

func (c *CreateTable) Down(ctx Context) ([]Stmt, error) {
	return nil, nil
}
