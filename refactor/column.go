package refactor

import (
	"encoding/xml"
	"fmt"

	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/refactor/strval"
)

type Column struct {
	XMLName xml.Name `xml:"Column"`

	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	IsNullable *string `xml:"isNullable,attr"`

	Default *string `xml:"default,attr"`

	Constraint *Constraint `xml:"Constraint"`
}

func (c *Column) Validate() error {
	if c.Name == "" {
		return ErrInvalid("Column.name cannot be empty")
	}
	if c.Type == "" {
		return ErrInvalid("Column.type cannot be empty")
	}
	if err := strval.ValidateBool(c.IsNullable); err != nil {
		return fmt.Errorf("Column.isNullable %v", err)
	}
	return c.Constraint.Validate()
}

func (c *Column) Definition(ctx Context) (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}

	expanded, err := ExpandAll(ctx, c.Name, c.Type)
	if err != nil {
		return "", err
	}
	name, t := expanded[0], expanded[1]

	result := fmt.Sprintf("%v %v", name, t)

	if !strval.Bool(c.IsNullable, false) {
		result = fmt.Sprintf("%v %v", result, dialect.NotNull)
	}
	if c.Default != nil {
		result = fmt.Sprintf("%v %v %v", result, dialect.Default, c.DefaultConstant(ctx))
	}
	return result, nil
}

func (c *Column) DefaultConstant(ctx Context) string {
	if c.Default == nil {
		return ""
	}
	return ctx.QuoteConst(*c.Default)
}
