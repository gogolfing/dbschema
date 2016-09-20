package refactor

import (
	"encoding/xml"
	"fmt"
)

const (
	errUnknownTokenType = xml.UnmarshalError("dbschema/refactor: unknown token type - xml syntax error")
)

type Validator interface {
	Validate() error
}

type ValidatorFunc func() error

func (vf ValidatorFunc) Validate() error {
	return vf()
}

func ValidateAll(validators ...Validator) error {
	for _, validator := range validators {
		if err := validator.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type ErrInvalid string

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("dbschema/refactor: invalid: %v", string(e))
}

func isXMLTokenEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}
