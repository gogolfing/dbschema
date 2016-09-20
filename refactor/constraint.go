package refactor

import "encoding/xml"

type Constraint struct {
	XMLName xml.Name `xml:"Constraint"`

	IsUnique   *BoolAttr   `xml:"isUnique,attr"`
	UniqueName *StringAttr `xml:"uniqueName,attr"`

	IsPrimaryKey *BoolAttr   `xml:"isPrimaryKey,attr"`
	PrimaryName  *StringAttr `xml:"primaryName,attr"`

	IsForeignKey *BoolAttr   `xml:"isForeignKey,attr"`
	ForeignName  *StringAttr `xml:"foreignName,attr"`
}

func (c *Constraint) Validate() error {
	return ValidateAll(
		c.IsUnique.Validator("Constraint.isUnique"),
		c.IsPrimaryKey.Validator("Constraint.isPrimaryKey"),
		c.IsForeignKey.Validator("Constraint.isForeignKey"),
	)
}
