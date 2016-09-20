package refactor

import (
	"encoding/xml"
	"fmt"

	"github.com/gogolfing/dbschema/dialect"
)

type Column struct {
	XMLName xml.Name `xml:"Column"`

	Name *StringAttr `xml:"name,attr"`
	Type *StringAttr `xml:"type,attr"`

	IsNullable *BoolAttr `xml:"isNullable,attr"`

	Default *StringAttr `xml:"default,attr"`

	Constraint *Constraint `xml:"Constraint"`
}

func (c *Column) Validate() error {
	return ValidateAll(
		c.Name.NotEmptyValidator("Column.name"),
		c.Type.NotEmptyValidator("Column.type"),
		c.IsNullable.Validator("Column.isNullable"),
		ValidatorFunc(func() error {
			if c.Constraint == nil {
				return nil
			}
			return c.Constraint.Validate()
		}),
	)
}

func (c *Column) Definition(ctx Context) (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}

	expanded, err := ExpandAll(
		ctx,
		c.Name,
		c.Type,
		c.IsNullable.Expander(true),
		c.Default,
	)
	if err != nil {
		return "", err
	}
	name, t, isNullable, def :=
		ctx.QuoteRef(expanded[0]),
		expanded[1],
		BoolString(expanded[2]),
		expanded[3]

	result := fmt.Sprintf("%v %v", name, t)

	if !isNullable {
		result = fmt.Sprintf("%v %v", result, dialect.NotNull)
	}
	if c.Default != nil {
		result = fmt.Sprintf("%v %v %v", result, dialect.Default, Constant(ctx, def, t))
	}

	return result, nil
}

func Constant(ctx Context, expr, t string) string {
	escaped, _ := ctx.EscapeConst(expr)
	return ctx.Cast(escaped, t)
}

func ColumnsValidator(columns []*Column) Validator {
	return ValidatorFunc(func() error {
		for _, c := range columns {
			if err := c.Validate(); err != nil {
				return err
			}
		}
		return nil
	})
}
