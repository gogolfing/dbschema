package dto

import "testing"

func TestNewString(t *testing.T) {
	value := "value"

	valuePtr := NewString(value)

	if *valuePtr != "value" {
		t.Fatal()
	}
}
