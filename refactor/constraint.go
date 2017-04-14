package refactor

type Constraint struct {
	IsUnique   NullBool
	UniqueName NullString

	IsPrimaryKey   NullBool
	PrimaryKeyName NullString
}
