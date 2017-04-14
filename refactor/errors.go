package refactor

import (
	"fmt"

	"github.com/gogolfing/dbschema/refactor/dto"
)

type UnknownDTOChangerError struct {
	dto.Changer
}

func (e *UnknownDTOChangerError) Error() string {
	return fmt.Sprintf("dbschema/refactor: unknown dto changer type %T", e.Changer)
}
