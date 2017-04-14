package refactor

type CreateTable struct {
	Name string

	IfNotExists NullBool

	Columns []*Column
}
