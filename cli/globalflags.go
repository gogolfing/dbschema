package cli

import "flag"

const (
	DefaultVerbose            = false
	DefaultConnectionFilePath = "connection.xml"
	DefaultChangeLogFilePath  = "changelog.xml"
	DefaultDBMS               = ""
	DefaultHost               = ""
	DefaultPort               = 0
	DefaultUser               = ""
	DefaultPassword           = ""
	DefaultDatabase           = ""
)

type globalFlags struct {
	verbose bool

	conn       string
	changeLog  string
	dbms       string
	host       string
	port       int
	user       string
	password   string
	database   string
	connParams connParams
}

func newGlobalFlags() *globalFlags {
	return &globalFlags{
		connParams: newConnParams(),
	}
}

func (gf *globalFlags) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&gf.verbose, "v", DefaultVerbose, "print verbose output")
	fs.StringVar(&gf.conn, "conn", DefaultConnectionFilePath, "path to connection file")
	fs.StringVar(&gf.changeLog, "changelog", DefaultChangeLogFilePath, "path to change log file")
	fs.StringVar(&gf.dbms, "dbms", DefaultDBMS, "the type of the dbms to connect to. this will override the value in -conn if not default")
	fs.StringVar(&gf.host, "host", DefaultHost, "host to connect to. this will override the value in -conn if not empty")
	fs.IntVar(&gf.port, "port", DefaultPort, "port to connect to. this will override the value in -conn if not empty")
	fs.StringVar(&gf.user, "user", DefaultUser, "user to connect as. this will override the value in -conn if not empty")
	fs.StringVar(&gf.password, "password", DefaultPassword, "password to connect with. this will override the value in -conn if not empty")
	fs.StringVar(&gf.database, "database", DefaultDatabase, "database to connect to. this will override the value in -conn if not empty")
	fs.Var(&gf.connParams, "conn-param", "list of connection parameters in the form of <name>=<value>. should be set with multiple flag definitions. these will override already set parameters in -conn")
}
