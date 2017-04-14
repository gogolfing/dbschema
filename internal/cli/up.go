package cli

import (
	"context"
	"io"

	"github.com/gogolfing/cli/subcommand"
)

type up struct{}

func newUpSubCommand() subcommand.SubCommand {
	return &subcommand.SubCommandStruct{
		NameValue:       "up",
		SynopsisValue:   "Applies up changes from the ChangeLog",
		ParameterSetter: &emptyParameterSetter{},
		ExecuteValue: func(ctx context.Context, _ io.Reader, out, outErr io.Writer) error {
			dbschema, logger, err := newDBSchemaLogger(ctx, out, outErr)
			if err != nil {
				return err
			}

			return execClose(dbschema, func() error {
				return dbschema.Up(logger, 0)
			})
		},
	}
}
