package dbschema

import (
	"fmt"
	"io"
)

func collectingAppliedChangeSets(w io.Writer) {
	fmt.Fprintln(w, "Collecting applied ChangeSets...\n")
}
