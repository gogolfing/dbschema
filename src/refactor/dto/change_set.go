package dto

import (
	"encoding/xml"
	"io"
	"os"
)

type Changer interface{}

type ChangeSet struct {
	XMLName xml.Name `xml:"ChangeSet"`

	Id string `xml:"id,attr"`

	Name *string `xml:"name,attr"`

	Author *string `xml:"author,attr"`

	Changers []Changer
}

func newChangeSet() *ChangeSet {
	return &ChangeSet{}
}

func NewChangeSetFile(path string) (*ChangeSet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	c, err := NewChangeSetReader(file)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewChangeSetReader(in io.Reader) (*ChangeSet, error) {
	dec := xml.NewDecoder(in)
	c := &ChangeSet{}
	err := dec.Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ChangeSet) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	c.XMLName = start.Name

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			c.Id = attr.Value
		case "name":
			c.Name = NewString(attr.Value)
		case "author":
			c.Author = NewString(attr.Value)
		}
	}

	changers, err := decodeInnerChangers(dec, start)
	if err != nil {
		return err
	}
	c.Changers = changers

	return nil
}

func decodeInnerChangers(dec *xml.Decoder, start xml.StartElement) ([]Changer, error) {
	result := []Changer{}
	for token, _ := dec.Token(); !IsXMLTokenEndElement(token); token, _ = dec.Token() {
		//we do not care about anything that is not an xml.StartElement.
		//and because of the for loop before, it cannot be an xml.EndElement.
		switch token.(type) {
		case xml.CharData, xml.Comment, xml.Directive, xml.ProcInst: //all other xml.Token types.
			continue
		}
		innerStart, ok := token.(xml.StartElement)
		if !ok {
			return nil, errUnknownTokenType
		}
		changer, err := decodeInnerChanger(dec, innerStart)
		if err != nil {
			return nil, err
		}
		if changer != nil {
			result = append(result, changer)
		}
	}
	return result, nil
}

func decodeInnerChanger(dec *xml.Decoder, innerStart xml.StartElement) (Changer, error) {
	var changer Changer

	switch innerStart.Name.Local {
	case "RawSql":
		changer = &RawSql{}
	default:
		return nil, UnknownChangerTypeError(innerStart.Name.Local)
	}

	if err := dec.DecodeElement(changer, &innerStart); err != nil {
		return nil, err
	}

	return changer, nil
}
