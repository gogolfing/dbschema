package dto

import "encoding/xml"

func newString(s string) *string {
	result := new(string)
	*result = s
	return result
}

func isXMLTokenEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}
