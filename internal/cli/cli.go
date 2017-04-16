package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gogolfing/cli"
	"github.com/gogolfing/cli/subcommand"
	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/dialect"
	"github.com/gogolfing/dbschema/logger"
	"github.com/gogolfing/dbschema/refactor"

	_ "github.com/lib/pq"
)

type key int

const keyGlobalFlags key = iota

const CommandName = "dbschema"

func Run(args []string) error {
	gf := &globalFlags{}

	sc := &subcommand.SubCommander{
		CommandName: CommandName,
		GlobalFlags: gf,
	}
	registerSubCommands(sc)

	ctx := withGlobalFlags(context.Background(), gf)

	fmt.Printf("%p\n", gf)

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
	fmt.Printf("%p\n", gf)

	logger := createLogger(gf, out, outErr)

	dialect, err := createDialect(conn)
	if err != nil {
		return nil, nil, &CreateDialectError{err}
	}
	changeLog, err := createChangeLog(gf)
	if err != nil {
		return nil, nil, &CreateChangeLogError{err}
	}
	dbschema, err := createDBSchema(dialect, conn, changeLog)
	if err != nil {
		return nil, nil, &CreateDBSchemaError{err}
	}
	return dbschema, logger, nil
}

func createLogger(gf *globalFlags, out, outErr io.Writer) logger.Logger {
	var verbose io.Writer
	if gf.verbose {
		verbose = out
	}
	return logger.NewLoggerWriters(verbose, out, out, outErr)
}

func createDialect(dbms string) (dialect.Dialect, error) {
	return dialect.NewDialect(dbms)
}

func createChangeLog(gf *globalFlags) (*refactor.ChangeLog, error) {
	return refactor.NewChangeLogFile(gf.changeLogPath)
}

func createDBSchema(d dialect.Dialect, conn *conn.Connection, cl *refactor.ChangeLog) (*dbschema.DBSchema, error) {
	return dbschema.Open(d, conn, cl)
}
