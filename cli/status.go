package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const statusShort = "prints status information"
const statusLong = `
status prints the status information about the current state of where the
database ChangeLog is.
This includes the current ChangeSet name, id, and author as well as how this
relates to the entire ChangeLog.
`

func createStatusSubCommand() SubCommand {
	return &subCommandStruct{
		FlagSetter: &statusFlags{
			&flagSetterStruct{
				name: SubCommandStatus,
			},
		},
		short:         statusShort,
		long:          statusLong,
		needsDBSchema: true,
		executor:      status,
	}
}

type statusFlags struct {
	*flagSetterStruct
}

func (sf *statusFlags) Usage(out io.Writer, _ *flag.FlagSet) {
	name := sf.Name()
	fmt.Fprintf(out, "Usage of %v: %v\n", name, name)
	fmt.Fprintf(
		out,
		"There are no %v or %v for %v\n",
		commandParametersName,
		commandOptionsName,
		name,
	)
}

func status(dbschema *dbschema.DBSchema, logger logger.Logger) error {
	err := fmt.Errorf("unimplemented")
	return err
}
