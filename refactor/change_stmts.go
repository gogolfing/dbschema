package refactor

import "encoding/xml"

type Changer interface {
	Stmts(ctx Context) (stmts []Stmt, err error)
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
		if changer != nil {
			result = append(result, changer)
		}
	}
	return result, nil
}

func decodeInnerChanger(dec *xml.Decoder, innerStart xml.StartElement) (Changer, error) {
	var changeStmts Changer = nil
	switch innerStart.Name.Local {
	case "AddPrimaryKey":
		changeStmts = &AddPrimaryKey{}
	case "CreateTable":
		changeStmts = &CreateTable{}
	case "RawSql":
		changeStmts = &RawSql{}
	default:
		return nil, dec.Skip()
	}
	if err := dec.DecodeElement(changeStmts, &innerStart); err != nil {
		return nil, err
	}
	return changeStmts, nil
}
