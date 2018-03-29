package dto

import (
	"encoding/xml"
	"io"
	"os"
	"path"
)

type ChangeLog struct {
	XMLName xml.Name `xml:"ChangeLog"`

	TableName     *string `xml:"tableName,attr"`
	LockTableName *string `xml:"lockTableName,attr"`

	Variables *Variables `xml:"Variables"`

	ChangeSets []*ChangeSet `xml:"ChangeSet"`

	path string
}

func NewChangeLogFile(path string) (*ChangeLog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewChangeLogReader(path, file)
}

func NewChangeLogReader(path string, in io.Reader) (*ChangeLog, error) {
	dec := xml.NewDecoder(in)
	c := &ChangeLog{}
	c.path = path
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ChangeLog) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	c.XMLName = start.Name
	c.ensureVariablesExist()
	c.copyAttributes(start)

	for token, _ := dec.Token(); !isXMLTokenEndElement(token); token, _ = dec.Token() {
		//we do not care about anything that is not an xml.StartElement.
		//and because of the for loop before, it cannot be an xml.EndElement.
		switch token.(type) {
		case xml.CharData, xml.Comment, xml.Directive, xml.ProcInst: //all other xml.Token types.
			continue
		}
		innerStart, ok := token.(xml.StartElement)
		if !ok {
			return errUnknownTokenType
		}
		if err := c.unmarshalXMLInnerElement(dec, innerStart); err != nil {
			return err
		}
	}

	return nil
}

func (c *ChangeLog) copyAttributes(start xml.StartElement) {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "tableName":
			c.TableName = newString(attr.Value)
		case "lockTableName":
			c.LockTableName = newString(attr.Value)
		}
	}
}

func (c *ChangeLog) unmarshalXMLInnerElement(dec *xml.Decoder, innerStart xml.StartElement) error {
	switch innerStart.Name.Local {
	case "Variables":
		return c.unmarshalXMLInnerVariables(dec, innerStart)
	case "Import":
		return c.unmarshalXMLInnerImport(dec, innerStart)
	case "ChangeSet":
		return c.unmarshalXMLInnerChangeSet(dec, innerStart)
	}
	return dec.Skip()
}

func (c *ChangeLog) unmarshalXMLInnerVariables(dec *xml.Decoder, startVariables xml.StartElement) error {
	vars := &Variables{}
	if err := dec.DecodeElement(vars, &startVariables); err != nil {
		return err
	}
	c.Variables.Values = append(c.Variables.Values, vars.Values...)
	return nil
}

func (c *ChangeLog) unmarshalXMLInnerImport(dec *xml.Decoder, startImport xml.StartElement) error {
	imp := &Import{}
	if err := dec.DecodeElement(imp, &startImport); err != nil {
		return err
	}
	if imp.Path == "" {
		return InvalidImportPathError(imp.Path)
	}
	path := c.importPath(imp.Path)
	cs, err := NewChangeSetFile(path)
	if err != nil {
		return err
	}
	c.ChangeSets = append(c.ChangeSets, cs)
	return nil
}

func (c *ChangeLog) unmarshalXMLInnerChangeSet(dec *xml.Decoder, startChangeSet xml.StartElement) error {
	cs := &ChangeSet{}
	if err := dec.DecodeElement(cs, &startChangeSet); err != nil {
		return err
	}
	c.ChangeSets = append(c.ChangeSets, cs)
	return nil
}

func (c *ChangeLog) importPath(relPath string) string {
	if c.path == "" {
		return relPath
	}
	return path.Join(path.Dir(c.path), relPath)
}

func (c *ChangeLog) ensureVariablesExist() {
	if c.Variables == nil {
		c.Variables = &Variables{}
	}
}

type Import struct {
	XMLName xml.Name `xml:"Import"`

	Path string `xml:"path,attr"`
}
