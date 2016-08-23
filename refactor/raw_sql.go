package refactor

import "encoding/xml"

type RawSql struct {
	XMLName xml.Name `xml:"RawSql"`

	Up   *string `xml:"Up"`
	Down *string `xml:"Down"`
}
