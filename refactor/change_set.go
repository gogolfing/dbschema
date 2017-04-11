package refactor

import (
	"encoding/xml"
	"io"
	"os"
)

const (
	errChangeSetCannotBeEmpty = ErrInvalid("ChangeSet cannot be empty")
)

type ChangeSet struct {
	XMLName xml.Name `xml:"ChangeSet"`

	Id string `xml:"id,attr"`

	Name *string `xml:"name,attr"`

	Author *string `xml:"author,attr"`

	changers []Changer
}

//add validation for an empty changeset.

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

func (c *ChangeSet) Stmts(ctx Context) (stmts []*Stmt, err error) {
	for _, changer := range c.changers {
		temp, err := changer.Stmts(ctx)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, temp...)
	}
	return
}

func (c *ChangeSet) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	c.XMLName = start.Name
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
	c.changers = changers
	return nil
}
