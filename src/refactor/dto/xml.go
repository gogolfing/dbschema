package dto

import "encoding/xml"

//IsXMLTokenEndElement returns whether or not token is an xml.EndElement.
func IsXMLTokenEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}
