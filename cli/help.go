package cli

import (
	"flag"
	"fmt"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

var errHelpExecution = fmt.Errorf("error executing help")

const HelpName = "help"

type help struct {
	subCommandName *string
}

func (h *help) Name() string {
	return HelpName
}

func (h *help) Aliases() []string {
	return nil
}

func (h *help) Synopsis() string {
	return fmt.Sprintf("Prints help information for a %v", subCommandName)
}

func (h *help) Usage() (params, paramDesc string, hasOptions bool) {
	params = formatArgument(subCommandName)
	paramDesc = fmt.Sprintf(
		"%v is the %v to print help information for",
		formatArgument(subCommandName),
		subCommandName,
	)
	return
}

func (h *help) Description() string {
	return h.Synopsis()
}

func (h *help) SetParameter(param string) error {
	if h.subCommandName != nil {
		return errInvalidParameter
	}
	h.subCommandName = new(string)
	*h.subCommandName = param
	return nil
}

func (h *help) SetFlags(_ *flag.FlagSet) {}

func (h *help) ValidateParameters() error {
	if h.subCommandName == nil {
		return errInvalidParameter
	}
	return nil
}

func (h *help) NeedsDBSchema() bool {
	return false
}

func (h *help) Execute(_ *dbschema.DBSchema, logger logger.Logger) error {
	sc, err := getSubCommand(*h.subCommandName)
	if err != nil {
		fmt.Fprintf(logger.Info(), "%v\n\n", err)
		printSubCommandsUsage(logger.Info())
		return errHelpExecution
	}

	f := flag.NewFlagSet("", flag.ContinueOnError)
	sc.SetFlags(f)
	printSubCommandUsage(logger.Info(), sc, f)

	return nil
}
