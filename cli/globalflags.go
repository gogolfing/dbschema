package cli

import (
	"flag"
	"fmt"
	"io"
)

const (
	DefaultVerbose            = false
	DefaultConnectionFilePath = "connection.xml"
	DefaultDBMS               = ""
	DefaultHost               = ""
	DefaultPort               = 0
	DefaultUser               = ""
	DefaultPassword           = ""
)

const (
	globalOptionsName     = "global_options"
	subCommandName        = "sub-command"
	commandParametersName = "command_parameters"
	commandOptionsName    = "command_options"
)

func printCommandUsage(out io.Writer, command string) {
	fmt.Fprintf(
		out,
		"Usage: %v [%v...] <%v> [%v...] [%v...]\n\n",
		command,
		globalOptionsName,
		subCommandName,
		commandParametersName,
		commandOptionsName,
	)
}

type globalFlags struct {
	command string

	verbose bool

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

func (gf *globalFlags) name() string {
	return globalOptionsName
}

func (gf *globalFlags) canHaveExtraArgs() bool {
	return true
}

func (gf *globalFlags) preParseArgs(args []string) []string {
	if len(args) == 0 {
		return []string{"-h"}
	}
	return args
}

func (gf *globalFlags) usage(out io.Writer, fs *flag.FlagSet) {
	fmt.Fprintf(out, "%v:\n", globalOptionsName)
	fs.PrintDefaults()
}

func (gf *globalFlags) set(fs *flag.FlagSet) {
	fs.BoolVar(&gf.verbose, "v", DefaultVerbose, "print verbose output")
	fs.StringVar(&gf.conn, "conn", DefaultConnectionFilePath, "path to connection file")
	fs.StringVar(&gf.dbms, "dbms", DefaultDBMS, "the type of the dbms to connect to. this will override the value in -conn if not default")
	fs.StringVar(&gf.host, "host", DefaultHost, "host to connect to. this will override the value in -conn if not empty")
	fs.IntVar(&gf.port, "port", DefaultPort, "port to connect to. this will override the value in -conn if not empty")
	fs.StringVar(&gf.user, "user", DefaultUser, "user to connect as. this will override the value in -conn if not empty")
	fs.StringVar(&gf.password, "password", DefaultPassword, "password to connect with. this will override the value in -conn if not empty")
	fs.Var(&gf.connParams, "conn-param", "list of connection parameters in the form of <name>=<value>. should be set with multiple flag definitions. these will override already set parameters in -conn")
}
