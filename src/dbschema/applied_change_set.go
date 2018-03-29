package dbschema

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/gogolfing/dbschema/src/logger"
)

type AppliedChangeSet struct {
	Id            string
	Name          *string
	Author        *string
	ExecutedAt    time.Time
	UpdatedAt     time.Time
	OrderExecuted int
	Sha256Sum     string
	Tags          []string
	Version       string
}

const changeSetRowLineFormat = "%-16s%v\n"

func (acs *AppliedChangeSet) String() string {
	buffer := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buffer, "ChangeSet - %s\n", acs.Id)
	if acs.Name != nil {
		fmt.Fprintf(buffer, changeSetRowLineFormat, "Name:", *acs.Name)
	}
	if acs.Author != nil {
		fmt.Fprintf(buffer, changeSetRowLineFormat, "Author:", *acs.Author)
	}
	fmt.Fprintf(buffer, changeSetRowLineFormat, "Executed At:", acs.ExecutedAt.Format(DefaultTimeFormat))

	return strings.TrimRight(buffer.String(), "\n")
}

func (acs *AppliedChangeSet) StringVerbose() string {
	buffer := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buffer, changeSetRowLineFormat, "SHA-256 Sum:", acs.Sha256Sum)
	tags := make([]string, 0, len(acs.Tags))
	for _, tag := range acs.Tags {
		tags = append(tags, fmt.Sprintf("%q", tag))
	}
	tagOutput := "<none>"
	if len(tags) > 0 {
		tagOutput = strings.Join(tags, ", ")
	}
	fmt.Fprintf(buffer, changeSetRowLineFormat, "Tags:", tagOutput)

	return strings.TrimRight(buffer.String(), "\n")
}

func printAppliedChangeSet(logger logger.Logger, acs *AppliedChangeSet) {
	fmt.Fprintln(logger.Info(), acs)
	fmt.Fprintln(logger.Verbose(), acs.StringVerbose())
	fmt.Fprintln(logger.Info())
}
