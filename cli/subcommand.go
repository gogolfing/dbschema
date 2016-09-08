package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

const (
	SubCommandStatus  = "status"
	SubCommandVersion = "version"
)

type errUnknownOrUndefinedSubCommand string

func (e errUnknownOrUndefinedSubCommand) Error() string {
	return fmt.Sprintf("unknown or undefined sub-command %q", string(e))
}

type SubCommandFunc func(dbschema *dbschema.DBSchema, logger logger.Logger) error

type SubCommand interface {
	FlagSetter

	ShortDescription() string
	LongDescription() string

	NeedsDBSchema() bool
	Execute(dbschema *dbschema.DBSchema, logger logger.Logger) error
}

type subCommandStruct struct {
	FlagSetter

	short string
	long  string

	needsDBSchema bool
	executor      SubCommandFunc
}

func (s *subCommandStruct) ShortDescription() string {
	return s.short
}

func (s *subCommandStruct) LongDescription() string {
	return s.long
}

func (s *subCommandStruct) NeedsDBSchema() bool {
	return s.needsDBSchema
}

func (s *subCommandStruct) Execute(dbschema *dbschema.DBSchema, logger logger.Logger) error {
	return s.executor(dbschema, logger)
}

type SubCommandCreator func() SubCommand

var subCommandCreators map[string]SubCommandCreator

func init() {
	subCommandCreators = map[string]SubCommandCreator{}

	subCommandCreators[SubCommandStatus] = createStatusSubCommand
	subCommandCreators[SubCommandVersion] = createVersionSubCommand
}

func getSubCommand(name string) (SubCommand, error) {
	creator, ok := subCommandCreators[name]
	if !ok {
		return nil, errUnknownOrUndefinedSubCommand(name)
	}
	return creator(), nil
}

func printSubCommandsUsage(out io.Writer) {
	names := []string{}
	shorts := map[string]string{}
	for name, creator := range subCommandCreators {
		sc := creator()
		names = append(names, name)
		shorts[name] = sc.ShortDescription()
	}
	sort.Strings(names)
	fmt.Fprintf(out, "sub-commands:\n")
	for _, name := range names {
		fmt.Fprintf(out, "  %v%v%v\n", name, pad(12, name), shorts[name])
	}
}

func pad(count int, name string) string {
	count = count - len(name)
	result := ""
	for count > 0 {
		result = result + " "
		count--
	}
	return result
}
