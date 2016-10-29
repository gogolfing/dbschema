package refactor

import (
	"encoding/xml"
	"fmt"
)

type StringAttr string

func NewStringAttr(s string) *StringAttr {
	sa := StringAttr(s)
	return &sa
}

func (s *StringAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	fmt.Println(attr)
	*s = StringAttr(attr.Value)
	return nil
}

func (s *StringAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if s == nil {
		return xml.Attr{}, nil
	}
	return xml.Attr{
		Name:  name,
		Value: string(*s),
	}, nil
}

func (s *StringAttr) NotEmptyValidator(path string) Validator {
	return ValidatorFunc(func() error {
		if s == nil || *s == "" {
			return ErrInvalid(fmt.Sprintf("%v cannot be empty", path))
		}
		return nil
	})
}

func (s *StringAttr) ExpandDefault(def string) Expander {
	return ExpanderFunc(func(ctx Context) (string, error) {
		if s == nil || *s == "" {
			return def, nil
		}
		return ctx.Expand(string(*s))
	})
}

func (s *StringAttr) Expand(ctx Context) (string, error) {
	if s == nil {
		return "", nil
	}
	return ctx.Expand(string(*s))
}
