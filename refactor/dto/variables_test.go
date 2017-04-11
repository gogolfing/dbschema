package dto

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestVariables_UnmarshalXML(t *testing.T) {
	raw := `
		<Variables>
			<Variable name="foo">bar</Variable>
		</Variables>
	`

	dec := xml.NewDecoder(strings.NewReader(raw))

	vars := &Variables{}

	if err := dec.Decode(vars); err != nil {
		t.Fatal(err)
	}

	if len(vars.Values) != 1 {
		t.Fatal()
	}
	if vars.Values[0].Name != "foo" || vars.Values[0].Value != "bar" {
		t.Fatal()
	}
}
