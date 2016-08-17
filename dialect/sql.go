package dialect

func NewSqlDialect() *Dialect {
	return &Dialect{
		CreateTable: "CREATE TABLE",

		Int:  "int",
		UUID: "UUID",
	}
}
