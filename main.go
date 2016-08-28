package main

import (
	"os"

	"github.com/gogolfing/dbschema/cli"
	"github.com/gogolfing/dbschema/dialect"
)

const (
	ExitCodeUnknownError             = 1
	ExitCodeGlobalParseFailed        = 2
	ExitCodeCreatingConnectionFailed = 3
	ExitCodeUnsupportedDBMS          = 4
)

func main() {
	cli := cli.NewCLI(os.Stdout, os.Stderr)

	err := cli.Run(os.Args[1:])

	if err != nil {
		exitError(err)
		return
	}
}

func exitError(err error) {
	code := ExitCodeUnknownError
	switch {

	case err == cli.ErrGlobalFlagParsingFailed:
		code = ExitCodeGlobalParseFailed

	case err == cli.ErrCreatingConnectionFailed:
		code = ExitCodeCreatingConnectionFailed

	case dialect.ErrUnsupportedDBMS:
		code = ExitCodeUnsupportedDBMS

	}
	exit(code)
}

/*
func exitError(err error) {
	codes := map[error]int{
		cli.ErrGlobalFlagParsingFailed:  ExitCodeGlobalParseFailed,
		cli.ErrCreatingConnectionFailed: ExitCodeCreatingConnectionFailed,
		dialect.ErrUnsupportedDBMS:      ExitCodeUnsupportedDBMS,
	}
	code, ok := codes[err]
	if !ok {
		fmt.Fprintln(os.Stderr, err)
		code = ExitCodeUnknownError
	}
	exit(code)
}
*/

func exit(code int) {
	os.Exit(code)
}
