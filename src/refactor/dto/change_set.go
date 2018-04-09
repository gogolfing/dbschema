package dto

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"github.com/gogolfing/dbschema/src/refactor"
)

type ChangeSets []*ChangeSet

func (cs ChangeSets) RefactorType() []*refactor.ChangeSet {
	result := make([]*refactor.ChangeSet, len(cs))
	for i, v := range cs {
		result[i] = v.RefactorType()
	}
	return result
}

func UnmarshalChangeSetXMLPath(path string) (*ChangeSet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return UnmarshalChangeSetXMLReader(file)
}

func UnmarshalChangeSetXMLReader(r io.Reader) (*ChangeSet, error) {
	dec := xml.NewDecoder(r)
	cs := newChangeSet()

	err := dec.Decode(cs)
	return cs, err
}

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

func (cs *ChangeSet) RefactorType() *refactor.ChangeSet {
	return &refactor.ChangeSet{
		Id:     cs.Id,
		Name:   refactor.NewNullString(cs.Name),
		Author: refactor.NewNullString(cs.Author),
		//TODO actually implement tags
		Tags:     nil,
		Changers: Changers(cs.Changers).RefactorType(),
	}
}

func (c *ChangeSet) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	log.Println("ChangeSet.UnmarshalXML()")

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
	log.Println("ChangeSet.UnmarshalXML()", "c.Changers", c.Changers, len(c.Changers))

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
		log.Printf("decodeInnerChanger result %T %v", changer, changer)
		if changer != nil {
			result = append(result, changer)
		}
	}
	return result, nil
}

func decodeInnerChanger(dec *xml.Decoder, start xml.StartElement) (Changer, error) {
	var changer Changer

	log.Println("decodeInnerChanger", start.Name.Local)

	switch start.Name.Local {
	case "RawSQL":
		changer = &RawSQL{}
	default:
		return nil, UnknownChangerTypeError(start.Name.Local)
	}

	if err := dec.DecodeElement(changer, &start); err != nil {
		return nil, err
	}

	return changer, nil
}
