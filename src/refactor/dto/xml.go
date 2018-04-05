package dto

import "encoding/xml"

//isXMLTokenEndElement returns whether or not token is an xml.EndElement.
func isXMLTokenEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}
