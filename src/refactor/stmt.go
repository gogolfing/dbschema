package refactor

//Stmt is a SQL statement that holds the raw SQL to execute on the database
//along with the variadic arguments for the SQL statement.
type Stmt struct {
	//Raw is the raw SQL statement.
	Raw string

	//Args are the arguments defined by placeholders in Raw.
	Args []interface{}
}
