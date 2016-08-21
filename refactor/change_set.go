package refactor

import (
	"encoding/xml"
	"io"
	"os"
)

const (
	errChangeSetCannotBeEmpty = ErrInvalid("ChangeSet cannot be empty")
	errInvalidChangeSetInner  = xml.UnmarshalError("dbschema/refactor: ChangeSet inner elements must be a valid refactor type")
	errUnknownTokenType       = xml.UnmarshalError("dbschema/refactor: unknown token type")
)

type ChangeSet struct {
	XMLName xml.Name `xml:"ChangeSet"`

	Id string `xml:"id,attr"`

	Name *string `xml:"name,attr"`

	Author *string `xml:"author,attr"`

	changers []Changer
}

func NewChangeSetFile(path string) (*ChangeSet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := NewChangeSetReader(file)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
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
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			c.Id = attr.Value
		case "name":
			c.Name = new(string)
			*c.Name = attr.Value
		case "author":
			c.Author = new(string)
			*c.Author = attr.Value
		}
	}
	changers, err := decodeInnerChangers(dec, start)
	if err != nil {
		return err
	}
	if len(changers) == 0 {
		return errChangeSetCannotBeEmpty
	}
	c.changers = changers
	return nil
}

func decodeInnerChangers(dec *xml.Decoder, start xml.StartElement) ([]Changer, error) {
	result := []Changer{}
	for token, _ := dec.Token(); !isXMLTokenEndElement(token); token, _ = dec.Token() {
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
		result = append(result, changer)
	}
	return result, nil
}

func decodeInnerChanger(dec *xml.Decoder, innerStart xml.StartElement) (Changer, error) {
	var changer Changer = nil
	switch innerStart.Name.Local {
	case "RawSql":
		changer = &RawSql{}
	}
	if changer == nil {
		return nil, errInvalidChangeSetInner
	}
	return changer, dec.DecodeElement(changer, &innerStart)
}
