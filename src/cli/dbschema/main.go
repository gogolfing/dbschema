package main

import (
	"os"

	"github.com/gogolfing/dbschema/src/cli"

	_ "github.com/lib/pq"
)

func main() {
	err := cli.Run(os.Args[1:])

	if err != nil {
		os.Exit(1)
	}
}
