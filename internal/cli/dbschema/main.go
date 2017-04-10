package main

import (
	"os"

	"github.com/gogolfing/dbschema/internal/cli"
)

func main() {
	err := cli.Run(os.Args[1:])

	if err != nil {
		os.Exit(1)
	}
}
