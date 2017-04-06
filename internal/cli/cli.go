package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gogolfing/cli"
	"github.com/gogolfing/cli/subcommand"
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

func executeDBSchemaLogger(ctx context.Context, in io.Reader, out, outErr io.Writer) error {
	// logger := logger.NewLoggerWriters(
	// 	out,
	// 	out,
	// 	out,
	// 	outErr,
	// )
	return nil
}
