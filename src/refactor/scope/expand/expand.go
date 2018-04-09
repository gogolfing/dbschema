package expand

import (
	"fmt"
	"strings"
)

//ReferenceDoesNotExistError is an error denoting a reference in an Expand string
//does not have a value associtated with it as determined by the MappingFunc.
type ReferenceDoesNotExistError struct {
	//Reference is the full reference value found in the Expand string.
	Reference string
}

//Error is the error implementation.
func (e *ReferenceDoesNotExistError) Error() string {
	return fmt.Sprintf("expand: reference %q does not exist", e.Reference)
}

//IsReferenceDoesNotExistError returns whether or not err is of type *ReferenceDoesNotExistError.
func IsReferenceDoesNotExistError(err error) bool {
	_, ok := err.(*ReferenceDoesNotExistError)
	return ok
}

//InvalidExpandSyntaxError is an error denoting an invalid reference or escape
//value in an Expand string.
type InvalidExpandSyntaxError struct {
	//Value is the Expand value that is invalid.
	Value string

	//Index is the index within Value at which the error is located.
	Index int

	//Description is a description of the syntax error.
	Description string
}

//Error is the error implementation.
func (e *InvalidExpandSyntaxError) Error() string {
	return fmt.Sprintf("expand: invalid expansion syntax for %q at index %d: %s", e.Value, e.Index, e.Description)
}

//MappingFunc is a function that accepts an input reference name and returns
//its associated value and an indicator as to whether or not it exists.
type MappingFunc func(value string) (result string, exists bool)

//Expand returns each reference in value replaced with the value being referenced.
//If a reference does not have a value associtated with it (determined by the ok
//return value of mapping) then a ReferenceDoesNotExistError error is returned.
//If the syntax of value is bad, then an InvalidExpandSyntaxError error is returned.
//
//When successful, all references in value are replaced (via recursive calls to
//Expand) with values returned from mapping.
func Expand(value string, mapping MappingFunc) (string, error) {
	var builder strings.Builder
	builder.Grow(2 * len(value))

	i := 0
	for j := 0; j < len(value); j++ {
		if value[j] == '$' {
			if j+1 >= len(value) {
				//invalid $<eof> escape.
				return "", &InvalidExpandSyntaxError{
					Value:       value,
					Index:       j,
					Description: "invalid $ at end of input",
				}
			} else {
				builder.WriteString(value[i:j])

				if value[j+1] == '$' {
					//double $ turns into a single $.
					builder.WriteRune('$')
					j += 1
					i = j + 1
				} else if value[j+1] == '{' {
					//start of reference.
					name, w := readName(value[j+2:])
					if w < 0 {
						//unclosed reference.
						return "", &InvalidExpandSyntaxError{
							Value:       value,
							Index:       j,
							Description: "unclosed reference",
						}
					}
					if name == "" {
						return "", &InvalidExpandSyntaxError{
							Value:       value,
							Index:       j,
							Description: "empty reference",
						}
					}
					mapped, ok := mapping(name)
					if !ok {
						return "", &ReferenceDoesNotExistError{
							Reference: name,
						}
					}

					innerMapped, err := Expand(mapped, mapping)
					if err != nil {
						return "", err
					}
					builder.WriteString(innerMapped)

					j += w + 1
					i = j + 1
				} else {
					//invalid $<byte> escape.
					return "", &InvalidExpandSyntaxError{
						Value:       value,
						Index:       j,
						Description: "invalid $<byte> escape",
					}
				}
			}
		}
	}

	//avoid copying if there is nothing that was expanded.
	if builder.Len() == 0 {
		return value, nil
	}

	builder.WriteString(value[i:])
	return builder.String(), nil
}

//readName reads value until it encounters a '}' to close a reference that has
//already started.
//An int return of -1 indicates there is no '}' ever encountered.
func readName(value string) (string, int) {
	for i := 0; i < len(value); i++ {
		if value[i] == '}' {
			return value[:i], i + 1
		}
	}
	return "", -1
}
