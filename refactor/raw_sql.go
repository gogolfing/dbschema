package refactor

import "encoding/xml"

type RawSql struct {
	XMLName xml.Name `xml:"RawSql"`

	Value string `xml:",innerxml"`
}
