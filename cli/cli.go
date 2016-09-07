package cli

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/logger"
)

var (
	ErrGlobalFlagParsingFailed     = errors.New("dbschema/cli: parsing global flags failed")
	ErrCreatingConnectionFailed    = errors.New("dbschmea/cli: could not create connection")
	ErrSubCommandFlagParsingFailed = errors.New("dbschema/cli: parsing command flags failed")
)

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
	gf, subCommandArgs, err := c.parseGlobalFlags(args)
	if err != nil {
		return ErrGlobalFlagParsingFailed
	}
	fmt.Fprintln(ioutil.Discard, gf, subCommandArgs)

	sc, err := c.parseSubCommand(subCommandArgs)
	if err != nil {
		return ErrSubCommandFlagParsingFailed
	}
	fmt.Fprintln(ioutil.Discard, sc)

	logger := c.createLogger(gf)

	if !sc.NeedsDBSchema() {
		return c.runWithoutDBSChema(sc, logger)
	}

	/*
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
	*/

	return nil
}

func (c *CLI) parseGlobalFlags(args []string) (gf *globalFlags, subCommandArgs []string, err error) {
	defer func() {
		if err != nil {
			c.printCommandUsage()
			c.printlnError(err)
			printSubCommandsUsage(c.outErr)
		}
	}()
	gf = newGlobalFlags()
	subCommandArgs, err = parseFlagSetter(gf, args)
	return
}

func (c *CLI) parseSubCommand(args []string) (sc SubCommand, err error) {
	defer func() {
		if err != nil {
			c.printCommandUsage()
			fmt.Fprintf(c.outErr, "%v", err)
			if _, ok := err.(errUnknownOrUndefinedSubCommand); ok {
				fmt.Fprintln(c.outErr, "\n")
				printSubCommandsUsage(c.outErr)
			} else {
				fmt.Fprintf(c.outErr, "\n%v\n", strings.TrimLeft(sc.LongDescription(), "\n"))
			}
		}
	}()
	if len(args) < 1 {
		err = errUnknownOrUndefinedSubCommand("")
		return
	}
	sc, err = getSubCommand(args[0])
	if err != nil {
		return
	}
	_, err = parseFlagSetter(sc, args[1:])
	return
}

func (c *CLI) createLogger(gf *globalFlags) logger.Logger {
	var verbose io.Writer = nil
	if gf.verbose {
		verbose = c.out
	}
	return logger.NewLoggerWriters(verbose, c.out, c.out, c.outErr)
}

func (c *CLI) runWithoutDBSChema(sc SubCommand, logger logger.Logger) error {
	return sc.Execute(nil, logger)
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

func (c *CLI) printCommandUsage() {
	printCommandUsage(c.outErr, c.command)
}

func (c *CLI) printlnError(err error) {
	fmt.Fprintln(c.outErr, err)
}
