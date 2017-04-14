package dbschema

import (
	"fmt"

	"github.com/gogolfing/dbschema/logger"
)

func (d *DBSchema) Status(logger logger.Logger) error {
	collectingAppliedChangeSets(logger.Verbose())

	empty := true
	err := d.ensureAppliedChangeSets(false, func(acs *AppliedChangeSet) {
		empty = false
		printAppliedChangeSet(logger, acs)
	})
	if err != nil {
		return err
	}

	if empty {
		fmt.Fprintln(logger.Info(), "No applied ChangeSets.")
	}
	return nil
}
