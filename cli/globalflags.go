package cli

import (
	"flag"
	"fmt"
	"io"
)

const (
	DefaultVerbose            = false
	DefaultConnectionFilePath = "connection.xml"
	DefaultChangeLogFilePath  = "changelog.xml"
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
	verbose bool

	conn       string
	changeLog  string
	dbms       string
	host       string
	port       int
	user       string
	password   string
	connParams connParams
}

func newGlobalFlags() *globalFlags {
	return &globalFlags{
		connParams: newConnParams(),
	}
}

func (gf *globalFlags) Name() string {
	return globalOptionsName
}

func (gf *globalFlags) CanHaveExtraArgs() bool {
	return true
}

func (gf *globalFlags) PreParseArgs(args []string) []string {
	if len(args) == 0 {
		return []string{"-h"}
	}
	return args
}

func (gf *globalFlags) Usage(out io.Writer, fs *flag.FlagSet) {
	fmt.Fprintf(out, "%v:\n", globalOptionsName)
	fs.PrintDefaults()
}

func (gf *globalFlags) Set(fs *flag.FlagSet) {
	fs.BoolVar(&gf.verbose, "v", DefaultVerbose, "print verbose output")
	fs.StringVar(&gf.conn, "conn", DefaultConnectionFilePath, "path to connection file")
	fs.StringVar(&gf.changeLog, "changelog", DefaultChangeLogFilePath, "path to change log file")
	fs.StringVar(&gf.dbms, "dbms", DefaultDBMS, "the type of the dbms to connect to. this will override the value in -conn if not default")
	fs.StringVar(&gf.host, "host", DefaultHost, "host to connect to. this will override the value in -conn if not empty")
	fs.IntVar(&gf.port, "port", DefaultPort, "port to connect to. this will override the value in -conn if not empty")
	fs.StringVar(&gf.user, "user", DefaultUser, "user to connect as. this will override the value in -conn if not empty")
	fs.StringVar(&gf.password, "password", DefaultPassword, "password to connect with. this will override the value in -conn if not empty")
	fs.Var(&gf.connParams, "conn-param", "list of connection parameters in the form of <name>=<value>. should be set with multiple flag definitions. these will override already set parameters in -conn")
}
