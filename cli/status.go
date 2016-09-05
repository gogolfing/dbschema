package cli

import (
	"flag"
	"io/ioutil"
)

type statusFlags struct {
}

func newStatusFlags(args []string) (*statusFlags, *flag.FlagSet, error) {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}
	return &statusFlags{}, fs, nil
}
