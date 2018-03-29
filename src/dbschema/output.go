package dbschema

import (
	"fmt"
	"io"
	"time"
)

const DefaultTimeFormat = time.RFC1123Z

func collectingAppliedChangeSets(w io.Writer) {
	fmt.Fprintln(w, "Collecting applied ChangeSets...\n")
}
