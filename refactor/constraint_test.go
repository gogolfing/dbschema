package refactor

import (
	"encoding/xml"
	"testing"
)

func TestConstraint_UnmarshalXML_defaults(t *testing.T) {
	input := []byte(`<Constraint />`)
	c := &Constraint{}
	err := xml.Unmarshal(input, c)
	if err != nil {
		t.Error(err)
	}
	if !c.IsNullableBool() || c.IsUniqueBool() || c.IsPrimaryBool() || c.IsForeignBool() {
		t.Fail()
	}
}
