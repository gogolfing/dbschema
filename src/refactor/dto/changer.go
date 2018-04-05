package dto

import (
	"encoding/xml"
	"fmt"

	"github.com/gogolfing/dbschema/src/refactor"
)

type UnknownRefactorChangerTypeError struct {
	refactor.Changer
}

func (e *UnknownRefactorChangerTypeError) Error() string {
	return fmt.Sprintf("refactor/dto: unknown refactor.Changer type %T", e.Changer)
}

type Changer interface{}

func MarshalRefactorChangerXML(enc *xml.Encoder, rc refactor.Changer) error {
	c, err := TransformRefactorChanger(rc)
	if err != nil {
		return err
	}
	return enc.Encode(c)
}

func TransformRefactorChanger(rc refactor.Changer) (Changer, error) {
	var result Changer

	switch rcType := rc.(type) {
	case *refactor.RawSQL:
		result = TransformRawSQL(rcType)

	default:
		return nil, &UnknownRefactorChangerTypeError{rc}
	}

	return result, nil
}
