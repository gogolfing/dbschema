package refactor

//Context is the interface used by Changers to be able to do their work.
//The contextual information provided by the interface helps Changers understand
//and produce the correct statements for refactoring.
type Context interface {
	//Expand returns an input value that has been expanded with all references
	//replaced with their associated values.
	//
	//An error should be returned if there is a syntax error in the input string
	//or a reference does not exist.
	Expand(string) (string, error)
}
