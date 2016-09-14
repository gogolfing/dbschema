package refactor

import (
	"encoding/xml"
	"fmt"

	"github.com/gogolfing/dbschema/refactor/strval"
)

type Constraint struct {
	XMLName xml.Name `xml:"Constraint"`

	IsUnique   *string `xml:"isUnique,attr"`
	UniqueName *string `xml:"uniqueName,attr"`

	IsPrimary   *string `xml:"isPrimary,attr"`
	PrimaryName *string `xml:"primaryName,attr"`

	IsForeign   *string `xml:"isForeign,attr"`
	ForeignName *string `xml:"foreignName,attr"`
}

func (c *Constraint) Validate() error {
	if err := strval.ValidateBool(c.IsUnique); err != nil {
		return fmt.Errorf("Constraint.isUnique %v", err)
	}
	if err := strval.ValidateBool(c.IsPrimary); err != nil {
		return fmt.Errorf("Constraint.isPrimary %v", err)
	}
	if err := strval.ValidateBool(c.IsForeign); err != nil {
		return fmt.Errorf("Constraint.isForeign %v", err)
	}
	return nil
}
