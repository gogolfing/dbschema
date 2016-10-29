package refactor

import "testing"

func TestAddPrimaryKey_Validate_errorOnEmptyTable(t *testing.T) {
	apk := &AddPrimaryKey{}
	err := apk.Validate()
	testStringAttrEmptyError(t, err, "AddPrimaryKey.table")
}

func TestAddPrimaryKey_Validate_errorOnIndexColumns(t *testing.T) {
	apk := &AddPrimaryKey{
		Table: NewStringAttr("table"),
		IndexColumns: []*IndexColumn{
			&IndexColumn{},
		},
	}
	err := apk.Validate()
	testStringAttrEmptyError(t, err, "IndexColumn.name")
}

func TestAddPrimaryKey_Validate_success(t *testing.T) {
	apk := &AddPrimaryKey{
		Table: NewStringAttr("table"),
		IndexColumns: []*IndexColumn{
			&IndexColumn{
				Name: NewStringAttr("name"),
			},
		},
	}
	if err := apk.Validate(); err != nil {
		t.Fail()
	}
}

func TestAddPrimaryKey_Stmts_errorValidate(t *testing.T) {
	apk := &AddPrimaryKey{}
	if _, err := apk.Stmts(defaultTestContext); err == nil {
		t.Fail()
	}
}

func TestAddPrimaryKey_Stmts_success(t *testing.T) {
	tests := []struct {
		ctx    Context
		apk    *AddPrimaryKey
		result string
	}{
		{
			defaultTestContext,
			&AddPrimaryKey{
				Table: NewStringAttr("table"),
			},
			`ALTER TABLE "table" ADD CONSTRAINT "table_pkey" PRIMARY KEY ()`,
		},
		{
			defaultTestContext,
			&AddPrimaryKey{
				Table: NewStringAttr("table"),
				Name:  NewStringAttr("name"),
			},
			`ALTER TABLE "table" ADD CONSTRAINT "name" PRIMARY KEY ()`,
		},
		{
			newTestContext(defaultTestDialect, "name1", "value1", "name2", "value2", "name3", "value3", "name4", "value4"),
			&AddPrimaryKey{
				Table: NewStringAttr("${name1}"),
				Name:  NewStringAttr("${name2}"),
				IndexColumns: []*IndexColumn{
					&IndexColumn{Name: NewStringAttr("${name3}")},
					&IndexColumn{Name: NewStringAttr("${name4}")},
				},
			},
			`ALTER TABLE "value1" ADD CONSTRAINT "value2" PRIMARY KEY ("value3", "value4")`,
		},
	}
	for _, test := range tests {
		stmts, err := test.apk.Stmts(test.ctx)
		if err != nil {
			t.Errorf("test.apk.Stmts(test.ctx) error = %v WANT %v", err, nil)
		}
		if len(stmts) != 1 {
			t.Errorf("test.apk.Stmts(test.ctx) len(stmts) = %v WANT %v", len(stmts), 1)
		}
		stmt := string(stmts[0])
		if stmt != test.result {
			t.Errorf("test.apk.Stmts(test.ctx) stmt = %v WANT %v", stmt, test.result)
		}
	}
}
