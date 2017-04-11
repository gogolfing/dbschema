package dto

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestStmt_UnmarshalXML(t *testing.T) {
	raw := `<Stmt> this is some raw stuff...</Stmt>`

	dec := xml.NewDecoder(strings.NewReader(raw))

	stmt := &Stmt{}
	if err := dec.Decode(stmt); err != nil {
		t.Fatal(err)
	}

	if stmt.Raw != " this is some raw stuff..." {
		t.Fatal()
	}
}
