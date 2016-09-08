package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
)

type FlagSetter interface {
	Name() string
	CanHaveExtraArgs() bool
	PreParseArgs([]string) []string
	Usage(io.Writer, *flag.FlagSet)
	Set(*flag.FlagSet)
}

type flagSetterStruct struct {
	name  string
	extra bool
}

func (fss *flagSetterStruct) Name() string {
	return fss.name
}

func (fss *flagSetterStruct) CanHaveExtraArgs() bool {
	return fss.extra
}

func (fss *flagSetterStruct) PreParseArgs(args []string) []string {
	return args
}

func (fss *flagSetterStruct) Set(_ *flag.FlagSet) {}

func parseFlagSetter(setter FlagSetter, args []string) ([]string, error) {
	fs := flag.NewFlagSet(setter.Name(), flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	setter.Set(fs)

	out := bytes.NewBuffer([]byte{})
	args = setter.PreParseArgs(args)
	if err := fs.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(out, "%v\n\n", err)
		}
		fs.SetOutput(out)
		setter.Usage(out, fs)
	}

	if out.Len() != 0 {
		return fs.Args(), fmt.Errorf("%v", out.String())
	}
	return fs.Args(), nil
}
