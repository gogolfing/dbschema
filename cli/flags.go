package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

type flagSetter interface {
	name() string
	modifyArgs([]string) []string
	set(fs *flag.FlagSet)
	usage() string
}

func parseFlagSetter(setter flagSetter, args []string) ([]string, error) {
	fs := flag.NewFlagSet(setter.name(), flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	setter.set(fs)
	out := bytes.NewBuffer([]byte{})
	if err := fs.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(out, "%v\n\n", err)
		}
		fmt.Fprintf(out, "%v\n", setter.usage())
		fs.SetOutput(out)
		fs.PrintDefaults()
	}

	if out.Len() != 0 {
		return fs.Args(), fmt.Errorf("%v", out.String())
	}
	return fs.Args(), nil
}

const (
	DefaultConnectionFilePath = "connection.xml"
	DefaultDBMS               = ""
	DefaultHost               = ""
	DefaultPort               = 0
	DefaultUser               = ""
	DefaultPassword           = ""
)

const (
	globalOptionsName     = "global_options"
	commandParametersName = "command_parameters"
	commandOptionsName    = "command_options"
)

type globalFlags struct {
	command string

	conn       string
	dbms       string
	host       string
	port       int
	user       string
	password   string
	connParams connParams
}

func newGlobalFlags(command string) *globalFlags {
	return &globalFlags{
		command:    command,
		conn:       "",
		connParams: newConnParams(),
	}
}

func (f *globalFlags) name() string {
	return globalOptionsName
}

func (f *globalFlags) usage() string {
	return fmt.Sprintf("Usage: %v [%v...] <command> [%v...] [%v...]", f.command, globalOptionsName, commandParametersName, commandOptionsName)
}

func (f *globalFlags) modifyArgs(args []string) []string {
	return args
}

func (f *globalFlags) set(fs *flag.FlagSet) {
	fs.StringVar(&f.conn, "conn", DefaultConnectionFilePath, "path to connection file")
	fs.StringVar(&f.dbms, "dbms", DefaultDBMS, "the type of the dbms to connect to. this will override the value in -conn if not default")
	fs.StringVar(&f.host, "host", DefaultHost, "host to connect to. this will override the value in -conn if not empty")
	fs.IntVar(&f.port, "port", DefaultPort, "port to connect to. this will override the value in -conn if not empty")
	fs.StringVar(&f.user, "user", DefaultUser, "user to connect as. this will override the value in -conn if not empty")
	fs.StringVar(&f.password, "password", DefaultPassword, "password to connect with. this will override the value in -conn if not empty")
	fs.Var(&f.connParams, "conn-param", "list of connection parameters in the form of <name>=<value>. should be set with multiple flag definitions. these will override already set parameters in -conn")
}
