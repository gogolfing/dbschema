package refactor

import "encoding/xml"

type Constraint struct {
	XMLName xml.Name `xml:"Constraint"`

	IsNullable string `xml:"isNullable,attr"`

	IsUnique   string `xml:"isUnique,attr"`
	UniqueName string `xml:"uniqueName,attr"`

	IsPrimary   string `xml:"isPrimary,attr"`
	PrimaryName string `xml:"primaryName,attr"`

	IsForeign   string `xml:"isForeign,attr"`
	ForeignName string `xml:"foreignName,attr"`
}

func (c *Constraint) IsNullableBool() bool {
	return BoolDefault(c.IsNullable, true)
}

func (c *Constraint) IsUniqueBool() bool {
	return BoolDefault(c.IsUnique, false)
}

func (c *Constraint) IsPrimaryBool() bool {
	return BoolDefault(c.IsPrimary, false)
}

func (c *Constraint) IsForeignBool() bool {
	return BoolDefault(c.IsForeign, false)
}
