package refactor

import "testing"

func TestRawSql_UnmarshalXML(t *testing.T) {
	source := `
		<RawSql>
			this is some raw "sql"
		</RawSql>
	`

	rs := &RawSql{}
	decodeSourceIntoValue(t, rs, source)

	if rs.Value != "\n\t\t\tthis is some raw \"sql\"\n\t\t" {
		t.Fail()
	}
}
