package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
)

var (
	ErrGlobalFlagParsingFailed  = errors.New("dbschema/cli: parsing global flags failed")
	ErrCreatingConnectionFailed = errors.New("dbschmea/cli: could not create connection")
	ErrCommandFlagParsingFailed = errors.New("dbschema/cli: parsing command flags failed")
)

type ErrUnknownOrUndefinedCommand string

func (e ErrUnknownOrUndefinedCommand) Error() string {
	return fmt.Sprintf("dbschema/cli: unknown or undefined sub-command %q", string(e))
}

type CLI struct {
	command string

	out    io.Writer
	outErr io.Writer
}

func NewCLI(command string, out, outErr io.Writer) *CLI {
	return &CLI{
		command: command,
		out:     out,
		outErr:  outErr,
	}
}

func (c *CLI) Run(args []string) error {
	gf, commandArgs, err := c.parseGlobalFlags(args)
	if err != nil {
		c.printError(err)
		return ErrGlobalFlagParsingFailed
	}

	conn, err := c.createConnection(gf)
	if err != nil {
		c.printlnError(err)
		return ErrCreatingConnectionFailed
	}
	fmt.Println(conn)

	dialect, err := c.createDialect(conn)
	if err != nil {
		c.printlnError(err)
		return err
	}
	fmt.Println(dialect)

	scFunc, scFlags, err := c.parseSubCommand(commandArgs)
	if err != nil {
		c.printlnError(err)
		return ErrCommandFlagParsingFailed
	}
	fmt.Println(scFunc == nil, scFlags)

	return nil
}

func (c *CLI) parseGlobalFlags(args []string) (*globalFlags, []string, error) {
	gf := newGlobalFlags(c.command)
	commandArgs, err := parseFlagSetter(gf, args)
	if err != nil {
		return nil, nil, err
	}
	return gf, commandArgs, nil
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

func (c *CLI) createDialect(conn *conn.Connection) (dialect.Dialect, error) {
	dbms, err := conn.DBMSValue()
	if err != nil {
		return nil, err
	}
	return dialect.NewDialect(dbms)
}

func (c *CLI) parseSubCommand(args []string) (subCommandFunc, subCommandFlags, error) {
	if len(args) < 1 {
		return nil, nil, ErrUnknownOrUndefinedCommand("")
	}
	command := args[0]
	commandArgs := args[1:]
	return c.parseSubCommandFromName(command, commandArgs)
}

func (c *CLI) printError(err error) {
	fmt.Fprint(c.outErr, err)
}

func (c *CLI) printlnError(err error) {
	fmt.Fprintln(c.outErr, err)
}
