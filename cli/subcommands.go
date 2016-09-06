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

type subCommandFunc func(dbschema *dbschema.DBSchema, logger logger.Logger) error

type subCommand interface {
	flagSetter

	description() string
	long() string

	needsDBSchema() bool
	execute(dbschema *dbschema.DBSchema, logger logger.Logger) error
}

type subCommandStruct struct {
	flagSetter

	descriptionValue string
	longValue        string

	needsDBSchemaValue bool
	executor           subCommandFunc
}

func newSubCommand(fs flagSetter, description, long string, needsDBSchema bool, executor subCommandFunc) subCommand {
	return &subCommandStruct{
		flagSetter:         fs,
		descriptionValue:   description,
		longValue:          long,
		needsDBSchemaValue: needsDBSchema,
		executor:           executor,
	}
}

func (s *subCommandStruct) description() string {
	return s.descriptionValue
}

func (s *subCommandStruct) long() string {
	return s.longValue
}

func (s *subCommandStruct) needsDBSchema() bool {
	return s.needsDBSchemaValue
}

func (s *subCommandStruct) execute(dbschema *dbschema.DBSchema, logger logger.Logger) error {
	return s.executor(dbschema, logger)
}

type subCommandCreator func() subCommand

var subCommandCreators map[string]subCommandCreator

func init() {
	subCommandCreators = map[string]subCommandCreator{}

	subCommandCreators[SubCommandVersion] = createVersionSubCommand
}

func getSubCommand(name string) (subCommand, error) {
	creator, ok := subCommandCreators[name]
	if !ok {
		return nil, errUnknownOrUndefinedSubCommand(name)
	}
	return creator(), nil
}

func printSubCommandsUsage(out io.Writer) {
	subCommandNames := []string{}
	subCommandDescriptions := map[string]string{}
	for name, creator := range subCommandCreators {
		subCommand := creator()
		subCommandNames = append(subCommandNames, name)
		subCommandDescriptions[name] = subCommand.description()
	}
	sort.Strings(subCommandNames)
	fmt.Fprintf(out, "sub-commands:\n")
	for _, name := range subCommandNames {
		fmt.Fprintf(out, "  %v%v%v\n", name, pad(name), subCommandDescriptions[name])
	}
}

func pad(name string) string {
	count := 12 - len(name)
	result := ""
	for count > 0 {
		result = result + " "
		count--
	}
	return result
}
