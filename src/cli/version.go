package cli

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/gogolfing/cli/subcommand"
	"github.com/gogolfing/dbschema/src/dbschema"
)

func newVersionSubCommand() subcommand.SubCommand {
	return &subcommand.SubCommandStruct{
		NameValue:        "version",
		SynopsisValue:    "Prints version and build information",
		DescriptionValue: "Prints the current version and build information.",
		ParameterSetter:  &emptyParameterSetter{},
		ExecuteValue: func(_ context.Context, _ io.Reader, out, _ io.Writer) error {
			_, err := fmt.Fprintf(
				out,
				"Version: %s\nBuild:   %s %s\n",
				dbschema.Version,
				runtime.GOOS,
				runtime.GOARCH,
			)
			return err
		},
	}
}
