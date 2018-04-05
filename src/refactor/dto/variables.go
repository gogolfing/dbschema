package dto

import (
	"encoding/xml"

	"github.com/gogolfing/dbschema/src/refactor"
)

type Variables struct {
	XMLName xml.Name `xml:"Variables"`

	Values []*Variable `xml:"Variable"`
}

func newVariables() *Variables {
	return &Variables{}
}

func (v *Variables) RefactorType() *refactor.Variables {
	result := refactor.NewVariables()

	for _, value := range v.Values {
		result.Put(value.Name, value.Value)
	}

	return result
}

type Variable struct {
	XMLName xml.Name `xml:"Variable"`

	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}
