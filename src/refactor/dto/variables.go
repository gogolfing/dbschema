package dto

import "encoding/xml"

type Variables struct {
	XMLName xml.Name `xml:"Variables"`

	Values []*Variable `xml:"Variable"`
}

func newVariables() *Variables {
	return &Variables{}
}

type Variable struct {
	XMLName xml.Name `xml:"Variable"`

	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}
