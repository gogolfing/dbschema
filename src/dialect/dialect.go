package dialect

//SQL constructs that are the same accross all supported dialects.
const (
	CreateTable = "CREATE TABLE"
	AlterTable  = "ALTER TABLE"
	Add         = "ADD"
	Constraint  = "CONSTRAINT"
	PrimaryKey  = "PRIMARY KEY"
	NotNull     = "NOT NULL"
	Default     = "DEFAULT"
	IfNotExists = "IF NOT EXISTS"
)

type Dialect interface {
	DBMS() string

	Syntax

	Types
}

type DialectStruct struct {
	DBMSValue string

	Syntax

	Types
}

func (d *DialectStruct) DBMS() string {
	return d.DBMSValue
}
