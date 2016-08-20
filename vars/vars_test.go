package vars

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestVariables_Dereference(t *testing.T) {
	source := `
<Variables>
	<Variable name="nameOne" value="valueOne" />
	<Variable name="nameTwo" value="valueTwo" />
</Variables>
`

	dec := xml.NewDecoder(strings.NewReader(source))

	vars := &Variables{}

	err := dec.Decode(vars)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(vars, &Variables{
		XMLName: xml.Name{"", "Variables"},
		Values: []*Variable{
			newVariable("nameOne", "valueOne"),
			newVariable("nameTwo", "valueTwo"),
		},
	}) {
		t.Fail()
	}

	os.Setenv("envVarYes", "yes")
	os.Setenv("envVarNo", "")

	tests := []struct {
		name  string
		value string
		err   error
	}{
		{"{nameOne}", "valueOne", nil},
		{"{nameTwo}", "valueTwo", nil},
		{"{nameThree}", "", ErrDoesNotExist},
		{"nameFour", "", ErrInvalidReference},

		{"${envVarYes}", "yes", nil},
		{"${envVarNo}", "", ErrDoesNotExist},
	}
	for _, test := range tests {
		value, err := vars.Dereference(test.name)
		if value != test.value || err != test.err {
			t.Errorf("vars.Dereference(%v) = %v, %v WANT %v, %v", test.name, value, err, test.value, test.err)
		}
	}
}

func newVariable(name, value string) *Variable {
	return &Variable{
		XMLName: xml.Name{"", "Variable"},
		Name:    name,
		Value:   value,
	}
}
