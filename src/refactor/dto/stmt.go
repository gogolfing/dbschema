package dto

import "encoding/xml"

type Stmt struct {
	XMLName xml.Name `xml:"Stmt"`

	Raw string `xml:",innerxml"`
}
