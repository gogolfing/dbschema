package refactor

import (
	"fmt"
	"testing"
)

func TestNewVariables(t *testing.T) {
	vars := NewVariables()

	if vars == nil || vars.values == nil {
		t.Fatal()
	}
}

func TestVaraibles_Len_ReturnsZeroWhenEmpty(t *testing.T) {
	vars := NewVariables()

	if len := vars.Len(); len != 0 {
		t.Fatal(len)
	}
}

func TestVariables_Len_ReturnsThePositiveNumberOfValuesSet(t *testing.T) {
	vars := NewVariables()

	for i := 0; i < 12; i++ {
		vars.Put(fmt.Sprintf("%d", i), "")
	}

	if len := vars.Len(); len != 12 {
		t.Fatal(len)
	}
}

func TestVariablesCanBePutReturningOldValueAndGetThoseSetValues(t *testing.T) {
	vars := NewVariables()

	if value, ok := vars.GetOk("name"); value != "" || ok {
		t.Fatal(value, ok)
	}

	if putResult := vars.Put("name", "value"); putResult != "" {
		t.Fatal(putResult)
	}

	if value, ok := vars.GetOk("name"); value != "value" || !ok {
		t.Fatal(value, ok)
	}

	if putResult := vars.Put("name", "value2"); putResult != "value" {
		t.Fatal(putResult)
	}

	if value, ok := vars.GetOk("name"); value != "value2" || !ok {
		t.Fatal(value, ok)
	}
}
