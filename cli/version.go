package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const versionShort = "prints the version"
const versionLong = `
version prints the current version information.
It ignores all ` + globalOptionsName + `.`

func createVersionSubCommand() SubCommand {
	return &subCommandStruct{
		FlagSetter: &versionFlags{
			&flagSetterStruct{
				name: SubCommandVersion,
			},
		},
		short:         versionShort,
		long:          versionLong,
		needsDBSchema: false,
		executor:      version,
	}
}

type versionFlags struct {
	*flagSetterStruct
}

func (vf *versionFlags) Usage(out io.Writer, _ *flag.FlagSet) {
	name := vf.Name()
	fmt.Fprintf(out, "Usage of %v: %v\n", name, name)
	fmt.Fprintf(
		out,
		"There are no %v or %v for %v\n",
		commandParametersName,
		commandOptionsName,
		name,
	)
}

func version(_ *dbschema.DBSchema, logger logger.Logger) error {
	_, err := fmt.Fprintln(logger.Info(), dbschema.Version)
	return err
}
