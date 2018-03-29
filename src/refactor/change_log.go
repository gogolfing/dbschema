package refactor

import (
	"github.com/gogolfing/dbschema/src/vars"
)

//ChangeLog is the top level entity that DBSchema works with.
//It contains a slice of ChangeSets that DBSchema will make on the database.
type ChangeLog struct {
	//TableName is name of the table that keeps track of the changes made in
	//the ChangeLog.
	TableName NullString

	//LockTableName is the name of the table that is used to lock access to
	//making changes to a single instance at a time.
	LockTableName NullString

	//Variables are the variables set in a ChangeLog.
	Variables *vars.Variables

	//ChangeSets are the ChangeSet(s) to apply on the database.
	ChangeSets []*ChangeSet
}
