package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gogolfing/cli"
	"github.com/gogolfing/cli/subcommand"
	"github.com/gogolfing/dbschema/src/dbschema"
	"github.com/gogolfing/dbschema/src/dialect"
	"github.com/gogolfing/dbschema/src/logger"
	"github.com/gogolfing/dbschema/src/refactor"
)

type key int

const keyGlobalFlags key = iota

const CommandName = "dbschema"

const (
	ENV_VAR_CONN_STRING = "DBSCHEMA_CONN"
	ENV_VAR_DBMS        = "DBSCHEMA_DBMS"
)

func Run(args []string) error {
	gf := &globalFlags{}

	sc := &subcommand.SubCommander{
		CommandName: CommandName,
		GlobalFlags: gf,
	}
	registerSubCommands(sc)

	ctx := withGlobalFlags(context.Background(), gf)

	err := sc.ExecuteContext(
		ctx,
		args,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
	if subcommand.IsExecutionError(err) {
		fmt.Fprintln(os.Stderr, err)
	}
	return err
}

func registerSubCommands(sc *subcommand.SubCommander) {
	sc.RegisterHelp("help", "", "")
	sc.RegisterList("list", "", "")

	sc.Register(newVersionSubCommand())
	sc.Register(newStatusSubCommand())
	sc.Register(newUpSubCommand())
	sc.Register(newDownSubCommand())
}

func withGlobalFlags(ctx context.Context, gf *globalFlags) context.Context {
	return context.WithValue(ctx, keyGlobalFlags, gf)
}

func globalFlagsFrom(ctx context.Context) *globalFlags {
	return ctx.Value(keyGlobalFlags).(*globalFlags)
}

type emptyParameterSetter struct{}

func (ps *emptyParameterSetter) ParameterUsage() ([]*cli.Parameter, string) {
	return nil, "There are no " + subcommand.ParametersName + " for this " + subcommand.SubCommandName
}

func (ps *emptyParameterSetter) SetParameters(params []string) error {
	if len(params) > 0 {
		return cli.ErrTooManyParameters
	}
	return nil
}

func execClose(dbschema *dbschema.DBSchema, f func() error) (err error) {
	defer func() {
		closeErr := dbschema.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	err = f()
	return
}

func newDBSchemaLogger(ctx context.Context, out, outErr io.Writer) (*dbschema.DBSchema, logger.Logger, error) {
	gf := globalFlagsFrom(ctx)

	conn, dbms := getConnAndDMBS(gf)
	logger := createLogger(gf, out, outErr)

	dialect, err := createDialect(dbms)
	if err != nil {
		return nil, nil, &CreateDialectError{err}
	}
	changeLog, err := createChangeLog(gf.changeLogPath)
	if err != nil {
		return nil, nil, &CreateChangeLogError{err}
	}
	dbschema, err := createDBSchema(dialect, conn, changeLog)
	if err != nil {
		return nil, nil, &CreateDBSchemaError{err}
	}
	return dbschema, logger, nil
}

func getConnAndDMBS(gf *globalFlags) (string, string) {
	conn := os.Getenv(ENV_VAR_CONN_STRING)
	if gf.conn != "" {
		conn = gf.conn
	}

	dbms := os.Getenv(ENV_VAR_DBMS)
	if gf.dbms != "" {
		dbms = gf.dbms
	}

	return conn, dbms
}

func createLogger(gf *globalFlags, out, outErr io.Writer) logger.Logger {
	var verbose io.Writer
	if gf.verbose {
		verbose = out
	}
	return logger.NewLoggerWriters(verbose, out, out, outErr)
}

func createDialect(dbms string) (dialect.Dialect, error) {
	if dbms == "" {
		return nil, fmt.Errorf("dbschema: value for dbms cannot be empty. please set the %s environment variable or the -dbms flag", ENV_VAR_DBMS)
	}
	return dialect.NewDialect(dbms)
}

func createChangeLog(path string) (*refactor.ChangeLog, error) {
	return refactor.NewChangeLogFile(path)
}

func createDBSchema(d dialect.Dialect, conn string, cl *refactor.ChangeLog) (*dbschema.DBSchema, error) {
	if conn == "" {
		return nil, fmt.Errorf("dbschema: value for conn cannot be empty. please set the %s environment variable or the -conn flag", ENV_VAR_CONN_STRING)
	}
	return dbschema.Open(d, conn, cl)
}
