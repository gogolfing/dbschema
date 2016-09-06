package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const versionDescription = "prints the version"
const versionLong = `
version prints the current version information.
It ignores all ` + globalOptionsName + `.`

func createVersionSubCommand() subCommand {
	return newSubCommand(
		&versionFlags{},
		versionDescription,
		versionLong,
		false,
		version,
	)
}

type versionFlags struct{}

func (vf *versionFlags) name() string {
	return SubCommandVersion
}

func (vf *versionFlags) canHaveExtraArgs() bool {
	return false
}

func (vf *versionFlags) preParseArgs(args []string) []string {
	return args
}

func (vf *versionFlags) usage(out io.Writer, _ *flag.FlagSet) {
	name := vf.name()
	fmt.Fprintf(out, "Usage of %v: %v\n", name, name)
	fmt.Fprintf(
		out,
		"There are no %v or %v for %v\n",
		commandParametersName,
		commandOptionsName,
		name,
	)
}

func (vf *versionFlags) set(_ *flag.FlagSet) {}

func version(_ *dbschema.DBSchema, logger logger.Logger) error {
	_, err := fmt.Fprintln(logger.Info(), dbschema.Version)
	return err
}
