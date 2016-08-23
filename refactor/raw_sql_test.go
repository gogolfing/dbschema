package refactor

import "testing"

func TestRawSql_UnmarshalXML(t *testing.T) {
	source := `
		<RawSql>
			<Up>this is some raw up</Up>
			<Down>this is some raw down</Down>
		</RawSql>
	`

	rs := &RawSql{}
	decodeSourceIntoValue(t, rs, source)

	if *rs.Up != "this is some raw up" {
		t.Fail()
	}
	if *rs.Down != "this is some raw down" {
		t.Fail()
	}
}
