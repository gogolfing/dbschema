package refactor

type Column struct {
	Name string
	Type string

	IsNullable NullBool

	Default NullString

	Constraint *Constraint
}
