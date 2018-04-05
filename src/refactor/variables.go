package refactor

//Variables is a collection of name, value pairs.
//
//NewVariables should be used to create instances of this type.
//The zero value for this type is not a in a valid state.
type Variables struct {
	values map[string]string
}

//NewVariables returns a new Variables instance ready for use.
func NewVariables() *Variables {
	return &Variables{
		values: map[string]string{},
	}
}

//Len returns the number of name, value pairs in v.
func (v *Variables) Len() int {
	return len(v.values)
}

//Put associates value with name in v.
//
//It returns the previously set value for name, if present.
func (v *Variables) Put(name, value string) string {
	current := v.values[name]

	v.values[name] = value

	return current
}

//GetOk returns the value associated with name and a bool indicating the value
//actually exists in v.
func (v *Variables) GetOk(name string) (string, bool) {
	result, ok := v.values[name]
	return result, ok
}
