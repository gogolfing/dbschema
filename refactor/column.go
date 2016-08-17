package refactor

import "encoding/xml"

type Column struct {
	XMLName xml.Name `xml:"Column"`

	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	Constraints []*Constraint `xml:"constraints"`
}
