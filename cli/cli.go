package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const (
	globalOptionsName     = "global_options"
	subCommandName        = "sub_command"
	commandParametersName = "command_parameters"
	commandOptionsName    = "command_options"
)

var (
	ErrGlobalFlagParsingFailed     = fmt.Errorf("dbschema/cli: parsing global flags failed")
	ErrSubCommandFlagParsingFailed = fmt.Errorf("dbschema/cli: parsing command flags failed")

	ErrCreatingConnectionFailed = fmt.Errorf("dbschema/cli: could not create Connection")
	ErrCreatingChangeLogFailed  = fmt.Errorf("dbschema/cli: could not create ChangeLog")

	ErrOpeningDBSchemaFailed = fmt.Errorf("dbschema/cli: could not open DBSchema")

	ErrExecutingSubCommandFailed = fmt.Errorf("dbschema/cli: executing sub-command failed")
)

var errUnsuppliedSubCommand = fmt.Errorf("%v not supplied", subCommandName)

func Run(command string, args []string, out, outErr io.Writer) error {
	gf, subCommandArgs, err := parseGlobalFlags(command, args, outErr)
	if err != nil {
		return ErrGlobalFlagParsingFailed
	}

	sc, err := parseSubCommand(subCommandArgs, outErr)
	if err != nil {
		return ErrSubCommandFlagParsingFailed
	}

	logger := createLogger(gf, out, outErr)

	if !sc.NeedsDBSchema() {
		return runWithDBSchema(sc, nil, logger, outErr)
	}

	return nil
}

func parseGlobalFlags(command string, args []string, outErr io.Writer) (gf *globalFlags, subCommandArgs []string, err error) {
	f := flag.NewFlagSet(command, flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	defer func() {
		if err != nil {
			if err != flag.ErrHelp {
				fmt.Fprintln(outErr, err)
				fmt.Fprintln(outErr)
			}
			printCommandUsage(outErr, command)
			f.SetOutput(outErr)
			f.PrintDefaults()
			if err == flag.ErrHelp {
				fmt.Fprintln(outErr)
				printSubCommandsUsage(outErr)
			}
		}
	}()

	gf = newGlobalFlags()
	gf.SetFlags(f)

	if len(args) == 0 {
		args = []string{"-h"}
	}
	err = f.Parse(args)
	subCommandArgs = f.Args()
	return
}

func printCommandUsage(out io.Writer, command string) {
	fmt.Fprintf(
		out,
		"Usage: %v %v %v %v %v\n",
		command,
		formatOptionalArgument(globalOptionsName, true),
		formatArgument(subCommandName),
		formatOptionalArgument(commandParametersName, true),
		formatOptionalArgument(commandOptionsName, true),
	)
}

func parseSubCommand(subCommandArgs []string, outErr io.Writer) (sc subCommand, err error) {
	var f *flag.FlagSet = nil

	defer func() {
		if err != nil {
			if err != flag.ErrHelp {
				fmt.Fprintln(outErr, err)
				fmt.Fprintln(outErr)
			}
			if _, unknown := err.(errUnknownSubCommand); unknown || err == errUnsuppliedSubCommand {
				printSubCommandsUsage(outErr)
			} else {
				printSubCommandUsage(outErr, sc, f)
			}
		}
	}()

	if len(subCommandArgs) == 0 {
		err = errUnsuppliedSubCommand
		return
	}

	name := subCommandArgs[0]
	sc, err = getSubCommand(name)
	if err != nil {
		return
	}

	f = flag.NewFlagSet(sc.Name(), flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	sc.SetFlags(f)

	parseArgs := subCommandArgs[1:]
	for len(parseArgs) > 0 && err == nil {
		err = f.Parse(parseArgs)
		parseArgs = f.Args()
		if err == nil && len(parseArgs) > 0 {
			err = sc.SetParameter(parseArgs[0])
			parseArgs = parseArgs[1:]
		}
	}
	if vpErr := sc.ValidateParameters(); err == nil && vpErr != nil {
		err = vpErr
	}

	return
}

func createLogger(gf *globalFlags, out, outErr io.Writer) logger.Logger {
	var verbose io.Writer = nil
	if gf.verbose {
		verbose = out
	}
	return logger.NewLoggerWriters(verbose, out, out, outErr)
}

func runWithDBSchema(sc subCommand, dbschema *dbschema.DBSchema, logger logger.Logger, outErr io.Writer) error {
	err := sc.Execute(dbschema, logger)
	if err != nil {
		if err != errHelpExecution {
			printlnError(outErr, err)
		}
		return ErrExecutingSubCommandFailed
	}
	return nil
}

func printlnError(outErr io.Writer, err error) {
	fmt.Fprintln(outErr, err)
}

func formatArgument(arg string) string {
	return fmt.Sprintf("<%v>", arg)
}

func formatOptionalArgument(arg string, many bool) string {
	result := "[" + arg
	if many {
		result += "..."
	}
	return result + "]"
}

/*
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
*/
