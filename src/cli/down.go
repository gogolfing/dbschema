package cli

import (
	"context"
	"io"

	"github.com/gogolfing/cli/subcommand"
)

type down struct{}

func newDownSubCommand() subcommand.SubCommand {
	return &subcommand.SubCommandStruct{
		NameValue:       "down",
		SynopsisValue:   "Rolls back already applied ChangeSets",
		ParameterSetter: &emptyParameterSetter{},
		ExecuteValue: func(ctx context.Context, _ io.Reader, out, outErr io.Writer) error {
			dbschema, logger, err := newDBSchemaLogger(ctx, out, outErr)
			if err != nil {
				return err
			}

			return execClose(dbschema, func() error {
				return dbschema.Down(logger, 0)
			})
		},
	}
}
