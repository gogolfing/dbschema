package dto

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"

	"github.com/gogolfing/dbschema/src/refactor"
)

const (
	TableNameAttr     = "tableName"
	LockTableNameAttr = "lockTableName"
)

const (
	ElementTypeVariables = "Variables"
	ElementTypeImport    = "Import"
	ElementTypeChangeSet = "ChangeSet"
)

//UnmarshalChangeLogXMLPath returns an XML parsed ChangeLog from the file at path.
func UnmarshalChangeLogXMLPath(path string) (*ChangeLog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return UnmarshalChangeLogXMLReader(path, file)
}

//UnmarshalChangeLogXMLReader returns an XML parsed ChangeLog from r.
func UnmarshalChangeLogXMLReader(path string, r io.Reader) (*ChangeLog, error) {
	cl := newChangeLog(path)

	dec := xml.NewDecoder(r)
	if err := dec.Decode(cl); err != nil {
		return nil, err
	}

	return cl, nil
}

//ChangeLog is the (un)marshalable type representing a refactor.ChangeLog.
type ChangeLog struct {
	XMLName xml.Name `xml:"ChangeLog"`

	TableName     *string `xml:"tableName,attr"`
	LockTableName *string `xml:"lockTableName,attr"`

	Variables *Variables `xml:"Variables"`

	ChangeSets []*ChangeSet `xml:"ChangeSet"`

	path string
}

func newChangeLog(path string) *ChangeLog {
	return &ChangeLog{
		Variables: &Variables{},
		path:      path,
	}
}

//RefactorType returns a refactor.ChangeLog that is equivalent to cl.
func (cl *ChangeLog) RefactorType() *refactor.ChangeLog {
	return &refactor.ChangeLog{
		TableName:     refactor.NewNullString(cl.TableName),
		LockTableName: refactor.NewNullString(cl.LockTableName),
		Variables:     cl.Variables.RefactorType(),
		ChangeSets:    ChangeSets(cl.ChangeSets).RefactorType(),
	}
}

//UnmarshalXML is the XMLUnmarshaler implementation.
func (cl *ChangeLog) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	cl.XMLName = start.Name
	cl.copyAttributes(start)

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
		if err := cl.unmarshalXMLInnerElement(dec, innerStart); err != nil {
			return err
		}
	}

	return nil
}

func (cl *ChangeLog) copyAttributes(start xml.StartElement) {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case TableNameAttr:
			cl.TableName = NewString(attr.Value)
		case LockTableNameAttr:
			cl.LockTableName = NewString(attr.Value)
		}
	}
}

func (cl *ChangeLog) unmarshalXMLInnerElement(dec *xml.Decoder, innerStart xml.StartElement) error {
	switch innerStart.Name.Local {
	case ElementTypeVariables:
		return cl.unmarshalXMLInnerVariables(dec, innerStart)
	case ElementTypeImport:
		return cl.unmarshalXMLInnerImport(dec, innerStart)
	case ElementTypeChangeSet:
		return cl.unmarshalXMLInnerChangeSet(dec, innerStart)
	}

	return dec.Skip()
}

func (cl *ChangeLog) unmarshalXMLInnerVariables(dec *xml.Decoder, start xml.StartElement) error {
	vars := newVariables()
	if err := dec.DecodeElement(vars, &start); err != nil {
		return err
	}

	cl.Variables.Values = append(cl.Variables.Values, vars.Values...)
	return nil
}

func (cl *ChangeLog) unmarshalXMLInnerImport(dec *xml.Decoder, start xml.StartElement) error {
	imp := newImport()
	if err := dec.DecodeElement(imp, &start); err != nil {
		return err
	}

	if imp.Path == "" {
		return InvalidImportPathError(imp.Path)
	}
	if err := imp.Validate(); err != nil {
		return err
	}

	cs, err := UnmarshalChangeSetXMLPath(cl.importPath(imp.Path))
	if err != nil {
		return err
	}

	cl.ChangeSets = append(cl.ChangeSets, cs)
	return nil
}

func (cl *ChangeLog) unmarshalXMLInnerChangeSet(dec *xml.Decoder, start xml.StartElement) error {
	cs := newChangeSet()
	if err := dec.DecodeElement(cs, &start); err != nil {
		return err
	}

	cl.ChangeSets = append(cl.ChangeSets, cs)
	return nil
}

func (cl *ChangeLog) importPath(relPath string) string {
	if cl.path == "" {
		return relPath
	}

	return filepath.Join(filepath.Dir(cl.path), relPath)
}
