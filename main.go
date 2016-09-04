package main

import (
	"os"

	"github.com/gogolfing/dbschema/cli"
	"github.com/gogolfing/dbschema/dialect"
)

const (
	ExitCodeUnknownError = 1

	ExitCodeParsingGlobalFlagsFailed = 2
	ExitCodeCreatingConnectionFailed = 3
	ExitCodeUnsupportedDBMS          = 4
	ExitCodeParsingCommandFailed     = 5
)

var cliArgs = os.Args[1:]

func main() {
	mainFunc()
}

var mainFunc = func() {
	cli := cli.NewCLI(os.Stdout, os.Stderr)

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
		cli.ErrGlobalFlagParsingFailed:  ExitCodeParsingGlobalFlagsFailed,
		cli.ErrCreatingConnectionFailed: ExitCodeCreatingConnectionFailed,
		cli.ErrCommandFlagParsingFailed: ExitCodeParsingCommandFailed,
	}

	code, ok := errCodes[err]
	if ok {
		return code
	}

	if _, ok := err.(cli.ErrUnknownOrUndefinedCommand); ok {
		return ExitCodeParsingCommandFailed
	}
	if _, ok := err.(dialect.ErrUnsupportedDBMS); ok {
		return ExitCodeUnsupportedDBMS
	}

	return ExitCodeUnknownError
}
