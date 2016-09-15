package refactor

import (
	"reflect"
	"testing"
)

func TestRawSql_UnmarshalXML(t *testing.T) {
	source := `
		<RawSql>
			<Up>
				<Stmt>this is some raw up</Stmt>
			</Up>
			<Down>
				<Stmt>this is some raw down</Stmt>
			</Down>
		</RawSql>
	`

	r := &RawSql{}
	decodeSourceIntoValue(t, r, source)

	if !reflect.DeepEqual(r.UpStmts, []Stmt{"this is some raw up"}) {
		t.Fail()
	}
	if !reflect.DeepEqual(r.DownStmts, []Stmt{"this is some raw down"}) {
		t.Fail()
	}
}

func TestRawSql_Validate(t *testing.T) {
	tests := []struct {
		*RawSql
		err error
	}{
		{
			&RawSql{
				UpStmts: []Stmt{"up"},
			},
			errDownMustBeNonEmpty,
		},
		{
			&RawSql{
				DownStmts: []Stmt{"down"},
			},
			errUpMustBeNonEmpty,
		},
		{
			&RawSql{
				UpStmts:   []Stmt{"up"},
				DownStmts: []Stmt{"down"},
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

func TestRawSql_Up(t *testing.T) {
	r := &RawSql{UpStmts: []Stmt{"up"}, DownStmts: []Stmt{"down"}}
	testStmtsFunc(t, r.Up, defaultTestContext, nil, "up")
}

func TestRawSql_Down(t *testing.T) {
	r := &RawSql{UpStmts: []Stmt{"up"}, DownStmts: []Stmt{"down"}}
	testStmtsFunc(t, r.Down, defaultTestContext, nil, "down")
}
