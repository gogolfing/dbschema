package main

import (
	"os"

	"github.com/gogolfing/dbschema/src/cli"
	"github.com/gogolfing/dbschema/src/dbschema"

	_ "github.com/lib/pq"
)

func init() {
	dbschema.DriverPostgresql = "postgres"
}

func main() {
	err := cli.Run(os.Args[1:])

	if err != nil {
		os.Exit(1)
	}
}
