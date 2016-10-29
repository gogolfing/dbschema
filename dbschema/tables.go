package dbschema

import "github.com/gogolfing/dbschema/refactor"

func createChangeLogTable(name string) *refactor.CreateTable {
	return &refactor.CreateTable{
		Name:        refactor.NewStringAttr(name),
		IfNotExists: refactor.NewBoolAttr(refactor.True),
		Columns: []*refactor.Column{
			&refactor.Column{
				Name:       refactor.NewStringAttr("changeset_id"),
				Type:       refactor.NewStringAttr("${Dialect.VarChar256}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
			&refactor.Column{
				Name: refactor.NewStringAttr("changeset_name"),
				Type: refactor.NewStringAttr("${Dialect.VarChar256}"),
			},
			&refactor.Column{
				Name: refactor.NewStringAttr("changeset_author"),
				Type: refactor.NewStringAttr("${Dialect.VarChar256}"),
			},
			&refactor.Column{
				Name:       refactor.NewStringAttr("executed_at"),
				Type:       refactor.NewStringAttr("${Dialect.Timestamp}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
			&refactor.Column{
				Name:       refactor.NewStringAttr("updated_at"),
				Type:       refactor.NewStringAttr("${Dialect.Timestamp}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
			&refactor.Column{
				Name:       refactor.NewStringAttr("order_executed"),
				Type:       refactor.NewStringAttr("${Dialect.Int32}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
			&refactor.Column{
				Name:       refactor.NewStringAttr("sha256_sum"),
				Type:       refactor.NewStringAttr("${Dialect.Char32}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
			&refactor.Column{
				Name:       refactor.NewStringAttr("tags"),
				Type:       refactor.NewStringAttr("${Dialect.VarChar1024}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
			&refactor.Column{
				Name:       refactor.NewStringAttr("dbschema_version"),
				Type:       refactor.NewStringAttr("${Dialect.VarChar32}"),
				IsNullable: refactor.NewBoolAttr(refactor.False),
			},
		},
	}
}

func createChangeLogLockTable(name string) *refactor.CreateTable {
	return nil
}
