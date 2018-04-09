package main

import (
	"os"

	"github.com/gogolfing/dbschema/src/cli"
	"github.com/gogolfing/dbschema/src/dbschema"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func init() {
	dbschema.DriverPostgresql = "postgres"
	dbschema.DriverMysql = "mysql"
}

func main() {
	err := cli.Run(os.Args[1:])

	if err != nil {
		os.Exit(1)
	}
}
