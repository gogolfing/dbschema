package dbschema

import (
	"fmt"
	"log"

	"github.com/gogolfing/dbschema/src/logger"
)

func (d *DBSchema) Status(logger logger.Logger) error {
	log.Println("starting status")
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
