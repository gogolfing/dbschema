package dto

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestRawSql_UnmarshalXML(t *testing.T) {
	raw := `
		<RawSQL>
			<Up>
				<Stmt>a</Stmt>
			</Up>
			<Down>
				<Stmt>b</Stmt>
			</Down>
		</RawSQL>
	`

	dec := xml.NewDecoder(strings.NewReader(raw))

	rawSql := &RawSQL{}
	if err := dec.Decode(rawSql); err != nil {
		t.Fatal(err)
	}

	if len(rawSql.Up.Stmts) != 1 || rawSql.Up.Stmts[0].Raw != "a" {
		t.Fatal()
	}
	if len(rawSql.Down.Stmts) != 1 || rawSql.Down.Stmts[0].Raw != "b" {
		t.Fatal()
	}
}
