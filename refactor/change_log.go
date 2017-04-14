package refactor

import (
	"github.com/gogolfing/dbschema/refactor/dto"
	"github.com/gogolfing/dbschema/vars"
)

type ChangeLog struct {
	TableName     NullString
	LockTableName NullString

	Variables *vars.Variables

	ChangeSets []*ChangeSet
}

func NewChangeLogFile(path string) (*ChangeLog, error) {
	dtoCl, err := dto.NewChangeLogFile(path)
	if err != nil {
		return nil, err
	}
	return newChangeLogDto(dtoCl)
}
