package main

import (
	"os"

	"github.com/gogolfing/dbschema/cli"
	"github.com/gogolfing/dbschema/dialect"
)

const (
	ExitCodeUnknownError = 1

	ExitCodeParsingGlobalFlagsFailed = 2
	ExitCodeParsingSubCommandFailed  = 5

	ExitCodeCreatingConnectionFailed = 3
	ExitCodeCreatingChangeLogFailed  = 6

	ExitCodeUnsupportedDBMS = 4

	ExitCodeOpenDBSchemaFailed = 7
)

var cliArgs = os.Args[1:]

func main() {
	mainFunc()
}

var mainFunc = func() {
	cli := cli.NewCLI(os.Args[0], os.Stdout, os.Stderr)

	err := cli.Run(cliArgs)

	if err != nil {
		exitError(os.Exit, getCode, err)
		return
	}
}

var exitError = func(exit func(int), getCode func(error) int, err error) {
	exit(getCode(err))
}

var getCode = func(err error) int {
	errCodes := map[error]int{
		cli.ErrGlobalFlagParsingFailed:     ExitCodeParsingGlobalFlagsFailed,
		cli.ErrSubCommandFlagParsingFailed: ExitCodeParsingSubCommandFailed,

		cli.ErrCreatingConnectionFailed: ExitCodeCreatingConnectionFailed,
		cli.ErrCreatingChangeLogFailed:  ExitCodeCreatingChangeLogFailed,

		cli.ErrOpeningDBSchemaFailed: ExitCodeOpenDBSchemaFailed,
	}

	code, ok := errCodes[err]
	if ok {
		return code
	}

	if _, ok := err.(dialect.ErrUnsupportedDBMS); ok {
		return ExitCodeUnsupportedDBMS
	}

	return ExitCodeUnknownError
}
