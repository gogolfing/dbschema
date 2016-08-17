package refactor

import "testing"

func TestBoolDefault(t *testing.T) {
	tests := []struct {
		attr   string
		def    bool
		result bool
	}{
		{"", true, true},
		{"", false, false},

		{"not true", true, false},
		{"not true", false, false},

		{"false", true, false},
		{"false", false, false},

		{"true", true, true},
		{"true", false, true},
	}
	for _, test := range tests {
		result := BoolDefault(test.attr, test.def)
		if result != test.result {
			t.Fail()
		}
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError("this is a validation error")

	if err.Error() != "refactor: validation error: this is a validation error" {
		t.Fail()
	}
}
