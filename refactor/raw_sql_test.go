package refactor

import (
	"reflect"
	"testing"
)

func TestRawSql_UnmarshalXML(t *testing.T) {
	source := `
		<RawSql>
			<Stmt>this is some raw stmt</Stmt>
			<Stmt>two</Stmt>
		</RawSql>
	`

	r := &RawSql{}
	decodeSourceIntoValue(t, r, source)

	if !reflect.DeepEqual(r.StmtSlice, []Stmt{"this is some raw stmt", "two"}) {
		t.Fail()
	}
}

func TestRawSql_Validate(t *testing.T) {
	tests := []struct {
		*RawSql
		err error
	}{
		{
			&RawSql{},
			errMustBeNonEmpty,
		},
		{
			&RawSql{
				StmtSlice: []Stmt{"stmt"},
			},
			nil,
		},
	}
	for _, test := range tests {
		err := test.RawSql.Validate()
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("test.RawSql.Validate() = %v WANT %v", err, test.err)
		}
	}
}

func TestRawSql_Stmts(t *testing.T) {
	r := &RawSql{StmtSlice: []Stmt{"stmt"}}
	testStmtsFunc(t, r.Stmts, defaultTestContext, nil, "stmt")
}
