package cli

import (
	"context"
	"io"
	"strings"

	"github.com/gogolfing/cli/subcommand"
)

func newStatusSubCommand() subcommand.SubCommand {
	return &subcommand.SubCommandStruct{
		NameValue:     "status",
		AliasesValue:  []string{"s"},
		SynopsisValue: "Prints status information",
		DescriptionValue: strings.Join(
			[]string{
				"Prints the status information about the currect state of where the database ChangeLog is.",
				"This includes the current ChangeSet name, id, and author as well as how this relates to the entire ChangeLog.",
			}, " ",
		),
		ParameterSetter: &emptyParameterSetter{},
		ExecuteValue: func(ctx context.Context, _ io.Reader, out, outErr io.Writer) error {
			dbschema, logger, err := newDBSchemaLogger(ctx, out, outErr)
			if err != nil {
				return err
			}

			return execClose(dbschema, func() error {
				return dbschema.Status(logger)
			})
		},
	}
}
