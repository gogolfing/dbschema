package vars

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestVariables_Merge(t *testing.T) {
	v := &Variables{
		values: map[string]string{
			"a": "a",
			"e": "e",
		},
	}

	other := &Variables{
		values: map[string]string{
			"a": "b",
			"c": "d",
		},
	}

	v.Merge(other)

	if !reflect.DeepEqual(v, &Variables{
		values: map[string]string{
			"a": "b",
			"c": "d",
			"e": "e",
		},
	}) {
		t.Fail()
	}
}

func TestVariables_Put(t *testing.T) {
	other := &Variable{
		Name:  "name",
		Value: "value",
	}

	v := &Variables{}

	v.Put(other)

	if !reflect.DeepEqual(v, &Variables{
		values: map[string]string{
			"name": "value",
		},
	}) {
		t.Fail()
	}
}

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
		values: map[string]string{
			"nameOne": "valueOne",
			"nameTwo": "valueTwo",
		},
	}) {
		t.Error(vars)
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
