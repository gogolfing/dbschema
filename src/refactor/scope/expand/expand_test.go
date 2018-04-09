package expand

import "testing"

func TestReferenceDoesNotExistError_Error(t *testing.T) {
	err := &ReferenceDoesNotExistError{
		Reference: "name",
	}

	if err.Error() != `expand: reference "name" does not exist` {
		t.Fatal(err)
	}
}

func TestInvalidExpandSyntaxError_Error(t *testing.T) {
	err := &InvalidExpandSyntaxError{
		Value:       "something to expand",
		Index:       12,
		Description: "desc",
	}

	if err.Error() != `expand: invalid expansion syntax for "something to expand" at index 12: desc` {
		t.Fatal(err)
	}
}

func TestExpand_ReturnsAllCorrectExpansions(t *testing.T) {
	cases := []struct {
		value   string
		result  string
		mapping MappingFunc
	}{
		{"", "", nil},
		{"hello", "hello", nil},
		{"}}", "}}", nil},
		{"$$ foobar $$$$", "$ foobar $$", nil},
		{
			"${whoami}",
			"foobar",
			mapMapper(map[string]string{
				"whoami": "foobar",
			}),
		},
		{
			"Equation: ${one} ${plus} ${two} = ${three} !!",
			"Equation: 1 + 2 = 3 !!",
			mapMapper(map[string]string{
				"one":   "1",
				"plus":  "+",
				"two":   "2",
				"three": "3",
			}),
		},
		{
			"${a}${b}${c}",
			"ABC",
			mapMapper(map[string]string{
				"a": "A",
				"b": "B",
				"c": "C",
			}),
		},
		{
			"${test}}{}",
			"result}{}",
			mapMapper(map[string]string{
				"test": "result",
			}),
		},
		{
			"${value1}",
			">>final<<",
			mapMapper(map[string]string{
				"value1": ">${value2}<",
				"value2": ">final<",
			}),
		},
	}

	for i, tc := range cases {
		result, err := Expand(tc.value, tc.mapping)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if result != tc.result {
			t.Errorf("%d: %q != %q", i, result, tc.result)
		}
	}
}

func TestExpand_ReturnsSyntaxErrorForAnUnclosedExpansion(t *testing.T) {
	result, err := Expand("hello ${world", nil)
	if result != "" {
		t.Fatal(result)
	}

	assertSyntaxError(t, err, "unclosed reference")
}

func TestExpand_ReturnsSyntaxErrorForAnEmptyReference(t *testing.T) {
	result, err := Expand("hello ${}", nil)
	if result != "" {
		t.Fatal(result)
	}

	assertSyntaxError(t, err, "empty reference")
}

func TestExpand_ReturnsSyntaxErrorForARecursiveEmptyReference(t *testing.T) {
	result, err := Expand(
		"${value1}",
		mapMapper(map[string]string{
			"value1": "${}",
		}),
	)
	if result != "" {
		t.Fatal(result)
	}

	assertSyntaxError(t, err, "empty reference")
}

func TestExpand_ReturnsSyntaxErrorForADollarUnknownByteExpansion(t *testing.T) {
	result, err := Expand("hello $whoami}", nil)
	if result != "" {
		t.Fatal(result)
	}

	assertSyntaxError(t, err, "invalid $<byte> escape")
}

func TestExpand_ReturnsSyntaxErrorForADollarAtTheEndOfInput(t *testing.T) {
	result, err := Expand("hello $", nil)
	if result != "" {
		t.Fatal(result)
	}

	assertSyntaxError(t, err, "invalid $ at end of input")
}

func assertSyntaxError(t *testing.T, err error, desc string) {
	se, ok := err.(*InvalidExpandSyntaxError)
	if !ok {
		t.Fatalf("%T", err)
	}

	if se.Description != desc {
		t.Fatal(se.Description, "!=", desc)
	}
}

func TestExpand_ReturnsReferenceDoesNotExistError(t *testing.T) {
	result, err := Expand("${value}", mapMapper(nil))
	if result != "" {
		t.Fatal(result)
	}

	if _, ok := err.(*ReferenceDoesNotExistError); !ok {
		t.Fatal(err)
	}
}

func mapMapper(m map[string]string) MappingFunc {
	return func(v string) (string, bool) {
		result, ok := m[v]
		return result, ok
	}
}
