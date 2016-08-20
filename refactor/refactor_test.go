package refactor

// import "testing"

// func TestStringDefaultBool(t *testing.T) {
// 	tr, f, notTrue := "true", "false", "not true"

// 	tests := []struct {
// 		value  *string
// 		def    bool
// 		result bool
// 	}{
// 		{nil, true, true},
// 		{nil, false, false},

// 		{&notTrue, true, false},
// 		{&notTrue, false, false},

// 		{&f, true, false},
// 		{&f, false, false},

// 		{&tr, true, true},
// 		{&tr, false, true},
// 	}
// 	for _, test := range tests {
// 		result := StringDefaultBool(test.value, test.def)
// 		if result != test.result {
// 			t.Fail()
// 		}
// 	}
// }

// func TestErrInvalid_Error(t *testing.T) {
// 	err := ErrInvalid("this is a validation error")

// 	if err.Error() != "refactor: invalid: this is a validation error" {
// 		t.Fail()
// 	}
// }
