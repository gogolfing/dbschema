package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const VersionName = "version"

type version struct{}

func (v *version) Name() string {
	return VersionName
}

func (v *version) Aliases() []string {
	return nil
}

func (v *version) Synopsis() string {
	return "Prints version and build information"
}

func (v *version) Usage() (params, paramDesc string, hasOptions bool) {
	return
}

func (v *version) Description() string {
	return "Prints the current version and build information."
}

func (v *version) SetParameter(_ string) error {
	return errInvalidParameter
}

func (v *version) SetFlags(_ *flag.FlagSet) {}

func (v *version) ValidateParameters() error {
	return nil
}

func (v *version) NeedsDBSchema() bool {
	return false
}

func (v *version) Execute(_ *dbschema.DBSchema, logger logger.Logger) error {
	_, err := fmt.Fprintf(
		logger.Info(),
		"Version: %v\nBuild:   %v %v\n",
		dbschema.Version,
		runtime.GOOS,
		runtime.GOARCH,
	)
	return err
}
