package cli

import (
	"flag"
	"strings"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const StatusName = "status"

var statusAliases = []string{"s"}

type status struct{}

func (s *status) Name() string {
	return StatusName
}

func (s *status) Aliases() []string {
	return statusAliases
}

func (s *status) Synopsis() string {
	return "Prints status information"
}

func (s *status) Usage() (params, paramDesc string, hasOptions bool) {
	return
}

func (s *status) Description() string {
	return strings.Join([]string{
		"Prints the status information about the current state of where the database ChangeLog is.",
		"This includes the current ChangeSet name, id, and author as well as how this relates to the entire ChangeLog.",
	}, " ")
}

func (s *status) SetParameter(_ string) error {
	return errInvalidParameter
}

func (s *status) SetFlags(_ *flag.FlagSet) {}

func (s *status) ValidateParameters() error {
	return nil
}

func (s *status) NeedsDBSchema() bool {
	return true
}

func (s *status) Execute(dbschema *dbschema.DBSchema, logger logger.Logger) error {
	return dbschema.Status(logger)
}
