package dbschema

import "github.com/gogolfing/dbschema/refactor"

//Constants for change log column names.
const (
	ColumnChangeSetId     = "changeset_id"
	ColumnChangeSetName   = "changeset_name"
	ColumnChangeSetAuthor = "changeset_author"
	ColumnExecutedAt      = "executed_at"
	ColumnUpdatedAt       = "updated_at"
	ColumnOrderExecuted   = "order_executed"
	ColumnSha256Sum       = "sha256_sum"
	ColumnTags            = "tags"
	ColumnVersion         = "dbschema_version"
)

//Constants for change log lock column names.
const (
	ColumnLockId   = "lock_id"
	ColumnIsLocked = "is_locked"
	ColumnLockedAt = "locked_at"
	ColumnLockedBy = "locked_by"
)

func createChangeLogTable(name string) *refactor.Stmt {
	return &refactor.Stmt{
		Raw: `
			CREATE TABLE IF NOT EXISTS ` + name + ` (
				changeset_id VARCHAR(256) NOT NULL,
				changeset_name VARCHAR(256),
				changeset_author VARCHAR(256),
				executed_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				order_executed INTEGER NOT NULL,
				sha256_sum CHAR(64) NOT NULL,
				tags VARCHAR(1024),
				dbschema_version VARCHAR(32) NOT NULL
			)
		`,
	}
	// return &refactor.CreateTable{
	// 	Name:        refactor.NewStringAttr(name),
	// 	IfNotExists: refactor.NewBoolAttr(refactor.True),
	// 	Columns: []*refactor.Column{
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("changeset_id"),
	// 			Type:       refactor.NewStringAttr("${Dialect.VarChar256}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name: refactor.NewStringAttr("changeset_name"),
	// 			Type: refactor.NewStringAttr("${Dialect.VarChar256}"),
	// 		},
	// 		&refactor.Column{
	// 			Name: refactor.NewStringAttr("changeset_author"),
	// 			Type: refactor.NewStringAttr("${Dialect.VarChar256}"),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("executed_at"),
	// 			Type:       refactor.NewStringAttr("${Dialect.Timestamp}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("updated_at"),
	// 			Type:       refactor.NewStringAttr("${Dialect.Timestamp}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("order_executed"),
	// 			Type:       refactor.NewStringAttr("${Dialect.Int32}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("sha256_sum"),
	// 			Type:       refactor.NewStringAttr("${Dialect.Char32}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("tags"),
	// 			Type:       refactor.NewStringAttr("${Dialect.VarChar1024}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr("dbschema_version"),
	// 			Type:       refactor.NewStringAttr("${Dialect.VarChar32}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 	},
	// }
}

func createChangeLogLockTable(name string) *refactor.Stmt {
	return &refactor.Stmt{
		Raw: `
			CREATE TABLE IF NOT EXISTS ` + name + ` (
				lock_id VARCHAR(32) NOT NULL,
				is_locked INTEGER NOT NULL,
				locked_at TIMESTAMP NOT NULL,
				locked_by VARCHAR(256) NOT NULL
			)
		`,
	}
	// return &refactor.CreateTable{
	// 	Name:        refactor.NewStringAttr(name),
	// 	IfNotExists: refactor.NewBoolAttr(refactor.True),
	// 	Columns: []*refactor.Column{
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr(ColumnLockId),
	// 			Type:       refactor.NewStringAttr("${Dialect.VarChar32}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr(ColumnIsLocked),
	// 			Type:       refactor.NewStringAttr("${Dialect.Integer}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr(ColumnLockedAt),
	// 			Type:       refactor.NewStringAttr("${Dialect.Timestamp}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 		&refactor.Column{
	// 			Name:       refactor.NewStringAttr(ColumnLockedBy),
	// 			Type:       refactor.NewStringAttr("${Dialect.VarChar256}"),
	// 			IsNullable: refactor.NewBoolAttr(refactor.False),
	// 		},
	// 	},
	// }
}
