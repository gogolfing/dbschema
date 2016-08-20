package refactor

// import (
// 	"encoding/xml"
// 	"fmt"
// )

// type Constraint struct {
// 	XMLName xml.Name `xml:"Constraint"`

// 	IsUnique   *string `xml:"isUnique,attr"`
// 	UniqueName *string `xml:"uniqueName,attr"`

// 	IsPrimary   *string `xml:"isPrimary,attr"`
// 	PrimaryName *string `xml:"primaryName,attr"`

// 	IsForeign   *string `xml:"isForeign,attr"`
// 	ForeignName *string `xml:"foreignName,attr"`
// }

// func (c *Constraint) IsUniqueBool() bool {
// 	return StringDefaultBool(c.IsUnique, false)
// }

// func (c *Constraint) IsPrimaryBool() bool {
// 	return StringDefaultBool(c.IsPrimary, false)
// }

// func (c *Constraint) IsForeignBool() bool {
// 	return StringDefaultBool(c.IsForeign, false)
// }

// func (c *Constraint) Validate() error {
// 	if err := ValidateStringBool(c.IsUnique); err != nil {
// 		return fmt.Errorf("Constraint.isUnique %v", err)
// 	}
// 	if err := ValidateStringBool(c.IsPrimary); err != nil {
// 		return fmt.Errorf("Constraint.isPrimary %v", err)
// 	}
// 	if err := ValidateStringBool(c.IsForeign); err != nil {
// 		return fmt.Errorf("Constraint.isForeign %v", err)
// 	}
// 	return nil
// }
