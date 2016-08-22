package refactor

import (
	"encoding/xml"
	"io"
	"os"
	"path"

	"github.com/gogolfing/dbschema/vars"
)

const (
	errInvalidImportPath     = ErrInvalid("Import.path cannot be empty")
	errInvalidChangeLogInner = xml.UnmarshalError("dbschema/refactor: ChangeLog inner types must be a Variables, Import, or ChangeSet")
)

type ChangeLog struct {
	XMLName xml.Name `xml:"ChangeLog"`

	path string

	Variables *vars.Variables `xml:"Variables"`

	ChangeSets []*ChangeSet `xml:"ChangeSet"`
}

func NewChangeLogFile(path string) (*ChangeLog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	c, err := NewChangeLogReader(path, file)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewChangeLogReader(path string, in io.Reader) (*ChangeLog, error) {
	dec := xml.NewDecoder(in)
	c := &ChangeLog{}
	c.path = path
	err := dec.Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ChangeLog) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	c.XMLName = start.Name
	c.ensureVariablesExist()
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

func (c *ChangeLog) unmarshalXMLInnerElement(dec *xml.Decoder, innerStart xml.StartElement) error {
	switch innerStart.Name.Local {
	case "Variables":
		return c.unmarshalXMLInnerVariables(dec, innerStart)
	case "Import":
		return c.unmarshalXMLInnerImport(dec, innerStart)
	case "ChangeSet":
		return c.unmarshalXMLInnerChangeSet(dec, innerStart)
	}
	return errInvalidChangeLogInner
}

func (c *ChangeLog) unmarshalXMLInnerVariables(dec *xml.Decoder, startVariables xml.StartElement) error {
	vars := &vars.Variables{}
	if err := dec.DecodeElement(vars, &startVariables); err != nil {
		return err
	}
	c.Variables.Merge(vars)
	return nil
}

func (c *ChangeLog) unmarshalXMLInnerImport(dec *xml.Decoder, startImport xml.StartElement) error {
	imp := &Import{}
	if err := dec.DecodeElement(imp, &startImport); err != nil {
		return err
	}
	if imp.Path == "" {
		return errInvalidImportPath
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
		c.Variables = &vars.Variables{}
	}
}

type Import struct {
	XMLName xml.Name `xml:"Import"`

	Path string `xml:"path,attr"`
}
