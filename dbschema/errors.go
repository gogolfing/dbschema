package dbschema

import "fmt"

type ParsingTagsError struct {
	Tags string
	Err  error
}

func (e *ParsingTagsError) Error() string {
	return fmt.Sprintf("dbschema: error parsing tags %q : %v", e.Tags, e.Err)
}
