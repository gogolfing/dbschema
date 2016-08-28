package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
)

const (
	DefaultConnectionFilePath = "connection.xml"
	DefaultDBMS               = ""
	DefaultHost               = "localhost"
	DefaultPort               = 0
	DefaultUser               = ""
	DefaultPassword           = ""
)

var (
	ErrGlobalFlagParsingFailed  = errors.New("dbschema/cli: parsing global flags failed")
	ErrCreatingConnectionFailed = errors.New("dbschmea/cli: could not create connection")
)

const (
	globalOptionsName     = "global_options"
	commandParametersName = "command_parameters"
	commandOptionsName    = "command_options"
)

type globalFlags struct {
	conn       string
	dbms       string
	host       string
	port       int
	user       string
	password   string
	connParams connParams
}

func newGlobalFlags() *globalFlags {
	return &globalFlags{
		conn:       "",
		connParams: newConnParams(),
	}
}

type CLI struct {
	out    io.Writer
	outErr io.Writer
}

func NewCLI(out, outErr io.Writer) *CLI {
	return &CLI{
		out:    out,
		outErr: outErr,
	}
}

func (c *CLI) Run(args []string) error {
	gf, commandArgs, err := c.parseGlobalFlags(args)
	if err != nil {
		return ErrGlobalFlagParsingFailed
	}

	conn, err := c.createConnection(gf)
	if err != nil {
		c.printlnError(err)
		return ErrCreatingConnectionFailed
	}

	dialect, err := dialect.NewDialect(conn.DBMS)
	if err != nil {
		c.printlnError(err)
		return err
	}
	fmt.Println(dialect)

	fmt.Println(commandArgs)
	return nil
}

func (c *CLI) parseGlobalFlags(osArgs []string) (*globalFlags, []string, error) {
	fs := flag.NewFlagSet(globalOptionsName, flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	fs.Usage = c.usage

	gf := &globalFlags{}
	fs.StringVar(&gf.conn, "conn", DefaultConnectionFilePath, "path to connection file")
	fs.StringVar(&gf.dbms, "dbms", DefaultDBMS, "the type of the dbms to connect to. this will override the value in -conn if not default")
	fs.StringVar(&gf.host, "host", DefaultHost, "host to connect to. this will override the value in -conn if not empty")
	fs.IntVar(&gf.port, "port", DefaultPort, "port to connect to. this will override the value in -conn if not empty")
	fs.StringVar(&gf.user, "user", DefaultUser, "user to connect as. this will override the value in -conn if not empty")
	fs.StringVar(&gf.password, "password", DefaultPassword, "password to connect with. this will override the value in -conn if not empty")
	fs.Var(&gf.connParams, "conn-param", "list of connection parameters in the form of <name>=<value>. should be set with multiple flag definitions. these will override already set parameters in -conn")

	err := fs.Parse(osArgs)
	fs.SetOutput(c.outErr)

	if err != nil {
		fs.PrintDefaults()
		return nil, nil, err
	}
	return gf, fs.Args(), err
}

func (c *CLI) createConnection(gf *globalFlags) (*conn.Connection, error) {
	conn, err := conn.NewConnectionFile(gf.conn)
	if err != nil {
		return nil, err
	}
	if gf.dbms != DefaultDBMS {
		conn.DBMS = gf.dbms
	}
	if gf.host != DefaultHost {
		conn.Host = gf.host
	}
	if gf.port != DefaultPort {
		conn.Port = gf.port
	}
	if gf.user != DefaultUser {
		conn.User = gf.user
	}
	if gf.password != DefaultPassword {
		conn.Password = gf.password
	}
	gf.connParams.eachParamValue(func(name, value string) {
		conn.PutParam(name, value)
	})
	return conn, nil
}

func (c *CLI) printlnError(err error) {
	fmt.Fprintln(c.outErr, err)
}

func (c *CLI) usage() {
	fmt.Fprintf(
		c.outErr,
		"Usage: %s [%s...] <command> [%s...] [%s...]\n",
		os.Args[0],
		globalOptionsName,
		commandParametersName,
		commandOptionsName,
	)
}
