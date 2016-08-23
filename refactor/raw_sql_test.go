package refactor

import (
	"reflect"
	"testing"
)

func TestRawSql_UnmarshalXML(t *testing.T) {
	source := `
		<RawSql>
			<Up>this is some raw up</Up>
			<Down>this is some raw down</Down>
		</RawSql>
	`

	rs := &RawSql{}
	decodeSourceIntoValue(t, rs, source)

	if *rs.UpValue != "this is some raw up" {
		t.Fail()
	}
	if *rs.DownValue != "this is some raw down" {
		t.Fail()
	}
}

func TestRawSql_Validate(t *testing.T) {
	ctx := newDefaultTestContext()
	tests := []struct {
		*RawSql
		err error
	}{
		{
			&RawSql{
				UpValue: newString("up"),
			},
			ErrInvalid("RawSql > Down must be defined"),
		},
		{
			&RawSql{
				DownValue: newString("down"),
			},
			ErrInvalid("RawSql > Up must be defined"),
		},
		{
			&RawSql{
				UpValue:   newString("up"),
				DownValue: newString("down"),
			},
			nil,
		},
	}
	for _, test := range tests {
		err := test.RawSql.Validate(ctx)
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("test.RawSql.Validate() = %v WANT %v", err, test.err)
		}
	}
}

func TestRawSql_Up(t *testing.T) {
	r := &RawSql{UpValue: newString("up")}
	testUpDownResult(t, r.Up, newDefaultTestContext(), nil, "up")
}

func TestRawSql_Down(t *testing.T) {
	r := &RawSql{DownValue: newString("down")}
	testUpDownResult(t, r.Down, newDefaultTestContext(), nil, "down")
}
