package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/logger"
	"github.com/gogolfing/dbschema/refactor"
)

var (
	ErrGlobalFlagParsingFailed     = fmt.Errorf("dbschema/cli: parsing global flags failed")
	ErrSubCommandFlagParsingFailed = fmt.Errorf("dbschema/cli: parsing command flags failed")

	ErrCreatingConnectionFailed = fmt.Errorf("dbschema/cli: could not create Connection")
	ErrCreatingChangeLogFailed  = fmt.Errorf("dbschema/cli: could not create ChangeLog")

	ErrOpeningDBSchemaFailed = fmt.Errorf("dbschema/cli: could not open DBSchema")

	ErrExecutingSubCommandFailed = fmt.Errorf("dbschema/cli: executing %v failed", subCommandName)
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

	sc, err := c.parseSubCommand(subCommandArgs)
	if err != nil {
		return ErrSubCommandFlagParsingFailed
	}

	logger := c.createLogger(gf)

	if !sc.NeedsDBSchema() {
		return c.runWithDBSChema(sc, nil, logger)
	}

	conn, err := c.createConnection(gf)
	if err != nil {
		c.printlnError(err)
		return ErrCreatingConnectionFailed
	}

	dialect, err := c.createDialect(conn)
	if err != nil {
		c.printlnError(err)
		return err
	}

	changeLog, err := c.createChangeLog(gf)
	if err != nil {
		c.printlnError(err)
		return ErrCreatingChangeLogFailed
	}

	dbschema, err := c.createDBSchema(dialect, conn, changeLog)
	if err != nil {
		c.printlnError(err)
		return ErrOpeningDBSchemaFailed
	}
	defer func() {
		if err := dbschema.Close(); err != nil {
			c.printlnError(err)
		}
	}()

	return c.runWithDBSChema(sc, dbschema, logger)
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
	if gf.database != DefaultDatabase {
		conn.Database = gf.database
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

func (c *CLI) createChangeLog(gf *globalFlags) (*refactor.ChangeLog, error) {
	return refactor.NewChangeLogFile(gf.changeLog)
}

func (c *CLI) createDBSchema(
	dialect dialect.Dialect,
	conn *conn.Connection,
	changeLog *refactor.ChangeLog,
) (*dbschema.DBSchema, error) {
	return dbschema.OpenSql(dialect, conn, changeLog)
}

func (c *CLI) runWithDBSChema(sc SubCommand, dbschema *dbschema.DBSchema, logger logger.Logger) error {
	err := sc.Execute(dbschema, logger)
	if err != nil {
		c.printlnError(err)
		return ErrExecutingSubCommandFailed
	}
	return nil
}

func (c *CLI) printCommandUsage() {
	printCommandUsage(c.outErr, c.command)
}

func (c *CLI) printlnError(err error) {
	fmt.Fprintln(c.outErr, err)
}
