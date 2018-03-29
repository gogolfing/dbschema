package refactor

//CreateTable is a Changer that creates a table on Up and drops a table on Down.
type CreateTable struct {
	//Name is the name of the table to create/drop.
	Name string

	//IfNotExists specifies to add "IF NOT EXISTS" to the statement.
	IfNotExists NullBool

	//Columns are the Column(s) to include in the create table statement.
	Columns []*Column
}

//Up is the Changer implementation.
//
//It returns Stmt(s) that create the table along with any keys or indexes required
//by the table.
func (ct *CreateTable) Up(ctx Context) ([]*Stmt, error) {
	return nil, nil
}

//Down is the Changer implementation.
//
//It returns Stmt(s) that drop all the resources created in Up in reverse created
//order
func (ct *CreateTable) Down(ctx Context) ([]*Stmt, error) {
	return nil, nil
}
