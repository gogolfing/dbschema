package dialect

import "testing"

var _ Dialect = &DialectStruct{}

func TestDialectStruct_DBMS(t *testing.T) {
	d := &DialectStruct{
		DBMSValue: "dbms",
	}

	if d.DBMS() != "dbms" {
		t.Fatal(d.DBMS())
	}
}
