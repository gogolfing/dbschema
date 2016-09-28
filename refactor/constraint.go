package refactor

import "encoding/xml"

type Constraint struct {
	XMLName xml.Name `xml:"Constraint"`

	IsUnique   *BoolAttr   `xml:"isUnique,attr"`
	UniqueName *StringAttr `xml:"uniqueName,attr"`

	IsPrimaryKey   *BoolAttr   `xml:"isPrimaryKey,attr"`
	PrimaryKeyName *StringAttr `xml:"primaryName,attr"`

	IsForeignKey   *BoolAttr   `xml:"isForeignKey,attr"`
	ForeignKeyName *StringAttr `xml:"foreignName,attr"`
}

func (c *Constraint) Validate() error {
	return ValidateAll(
		c.IsUnique.Validator("Constraint.isUnique"),
		c.IsPrimaryKey.Validator("Constraint.isPrimaryKey"),
		c.IsForeignKey.Validator("Constraint.isForeignKey"),
	)
}

func (c *Constraint) primaryKey(ctx Context) (string, bool, error) {
	expanded, err := ExpandAll(ctx, c.PrimaryKeyName, c.IsPrimaryKey.Expander(false))
	if err != nil {
		return "", false, err
	}
	return expanded[0], BoolString(expanded[1]), nil
}
