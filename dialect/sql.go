package dialect

func NewSqlDialect() *Dialect {
	return &Dialect{
		CreateTable: "CREATE TABLE",

		NotNull: "NOT NULL",
		Default: "DEFAULT",

		Int:  "int",
		UUID: "UUID",
	}
}
