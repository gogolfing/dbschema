package refactor

import (
	"encoding/xml"
	"fmt"
)

const (
	errUnknownTokenType = xml.UnmarshalError("dbschema/refactor: unknown token type - xml syntax error")
)

const (
	True  = "true"
	False = "false"
)

func StringDefault(value *string, def string) string {
	if value == nil {
		return def
	}
	return *value
}

func StringDefaultBool(value *string, def bool) bool {
	if value == nil {
		return def
	}
	return *value == True
}

func ValidateStringBool(value *string) error {
	if value == nil {
		return nil
	}
	if *value != True || *value != False {
		return fmt.Errorf("must be %q, %q, or not present", True, False)
	}
	return nil
}

type Validator interface {
	Validate(ctx Context) error
}

type ErrInvalid string

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("dbschema/refactor: invalid: %v", string(e))
}

func StmtsFromFuncs(ctx Context, stmtFuncs ...func(ctx Context) (stmt string, err error)) (stmts []string, err error) {
	for _, stmtFunc := range stmtFuncs {
		stmt, err := stmtFunc(ctx)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func isXMLTokenEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}
