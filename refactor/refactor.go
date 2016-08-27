package refactor

import (
	"encoding/xml"
	"fmt"
)

const (
	errUnknownTokenType = xml.UnmarshalError("dbschema/refactor: unknown token type - xml syntax error")
)

type Validator interface {
	Validate(ctx Context) error
}

type ErrInvalid string

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("dbschema/refactor: invalid: %v", string(e))
}

func isXMLTokenEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}
