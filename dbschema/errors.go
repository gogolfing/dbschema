package dbschema

import (
	"fmt"

	"github.com/gogolfing/dbschema/refactor"
)

type ErrChangeSetOutOfOrder struct {
	ChangeSets []*refactor.ChangeSet

	*ChangeSetRow

	//If not Before, then after.
	BeforeRow bool
}

func (e *ErrChangeSetOutOfOrder) Error() string {
	return fmt.Sprintf("%v", *e)
}
