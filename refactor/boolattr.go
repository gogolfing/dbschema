package refactor

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/gogolfing/dbschema/vars"
)

const (
	True  = "true"
	False = "false"
)

type BoolAttr string

func NewBoolAttr(s string) *BoolAttr {
	ba := BoolAttr(s)
	return &ba
}

func (b *BoolAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	*b = BoolAttr(attr.Value)

	if *b == "" {
		return nil
	}

	_, err := strconv.ParseBool(string(*b))
	if err != nil && !vars.IsVariableReference(string(*b)) {
		return err
	}
	return nil
}

func (b *BoolAttr) Validator(path string) Validator {
	return ValidatorFunc(func() error {
		if b == nil {
			return nil
		}
		_, err := strconv.ParseBool(string(*b))
		if err != nil && !vars.IsVariableReference(string(*b)) {
			return ErrInvalid(fmt.Sprintf("%v must be a boolean or variable expression", path))
		}
		return nil
	})
}

func (b *BoolAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if b == nil {
		return xml.Attr{}, nil
	}
	return xml.Attr{
		Name:  name,
		Value: string(*b),
	}, nil
}

func (b *BoolAttr) Expander(def bool) Expander {
	return ExpanderFunc(func(ctx Context) (string, error) {
		toParse := fmt.Sprint(def)
		if b != nil {
			toParse = string(*b)
		}
		value, err := ctx.Expand(toParse)
		if err != nil {
			return "", err
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return "", err
		}
		if boolValue {
			return True, nil
		}
		return False, nil
	})
}

func BoolString(value string) bool {
	return value == True
}
