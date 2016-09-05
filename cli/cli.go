package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dialect"
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
	gf, subCommandArgs, err := c.runGlobalFlags(args)
	if err != nil {
		return ErrGlobalFlagParsingFailed
	}
	fmt.Fprintln(ioutil.Discard, gf, subCommandArgs)

	subCommand, err := c.runSubCommand(subCommandArgs)
	if err != nil {
		return ErrSubCommandFlagParsingFailed
	}
	fmt.Fprintln(ioutil.Discard, subCommand)

	subCommand.execute(nil, nil)

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

func (c *CLI) runGlobalFlags(args []string) (gf *globalFlags, subCommandArgs []string, err error) {
	defer func() {
		if err != nil {
			c.printCommandUsage()
			c.printlnError(err)
			printSubCommandsUsage(c.outErr)
		}
	}()
	gf = newGlobalFlags(c.command)
	subCommandArgs, err = parseFlagSetter(gf, args)
	return
}

func (c *CLI) runSubCommand(args []string) (sc subCommand, err error) {
	defer func() {
		if err != nil {
			c.printCommandUsage()
			fmt.Fprintf(c.outErr, "%v", err)
			if _, ok := err.(errUnknownOrUndefinedSubCommand); ok {
				fmt.Fprintln(c.outErr, "\n")
				printSubCommandsUsage(c.outErr)
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

type flagSetter interface {
	name() string
	canHaveExtraArgs() bool
	preParseArgs([]string) []string
	usage(io.Writer, *flag.FlagSet)
	set(*flag.FlagSet)
}

func parseFlagSetter(setter flagSetter, args []string) ([]string, error) {
	fs := flag.NewFlagSet(setter.name(), flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	setter.set(fs)

	out := bytes.NewBuffer([]byte{})
	args = setter.preParseArgs(args)
	if err := fs.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(out, "%v\n\n", err)
		}
		fs.SetOutput(out)
		setter.usage(out, fs)
	}

	if out.Len() != 0 {
		return fs.Args(), fmt.Errorf("%v", out.String())
	}
	return fs.Args(), nil
}
