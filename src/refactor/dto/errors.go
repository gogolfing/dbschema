package dto

import (
	"encoding/xml"
	"fmt"
)

const errUnknownTokenType = xml.UnmarshalError("dbschema/refactor: unknown token type - xml syntax error")

type UnknownChangerTypeError string

func (e UnknownChangerTypeError) Error() string {
	return fmt.Sprintf("dbschema/refactor/dto: unknown Changer type in ChangeSet %q", string(e))
}
