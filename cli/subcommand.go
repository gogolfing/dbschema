package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gogolfing/dbschema/dbschema"
	"github.com/gogolfing/dbschema/logger"
)

var errInvalidParameter = fmt.Errorf("invalid %v", commandParametersName)

type errUnknownSubCommand string

func (e errUnknownSubCommand) Error() string {
	return fmt.Sprintf("unknown sub-command %q", string(e))
}

type subCommand interface {
	Name() string
	Aliases() []string
	Synopsis() string

	//Usage returns the command parameters, command parameters description, and
	//an indication of whether or not the subCommand has options.
	//params and paramDesc are written to output as is if applicable.
	Usage() (params, paramDesc string, hasOptions bool)

	Description() string

	SetParameter(param string) error
	SetFlags(f *flag.FlagSet)

	ValidateParameters() error

	NeedsDBSchema() bool
	Execute(dbschema *dbschema.DBSchema, logger logger.Logger) error
}

var subCommandRegistry = struct {
	names   map[string]subCommand
	aliases map[string]subCommand
}{
	names:   map[string]subCommand{},
	aliases: map[string]subCommand{},
}

func init() {
	registerSubCommands(
		&version{},
		&status{},
		&help{},
	)
}

func registerSubCommands(subCommands ...subCommand) {
	for _, sc := range subCommands {
		subCommandRegistry.names[sc.Name()] = sc
		for _, alias := range sc.Aliases() {
			subCommandRegistry.aliases[alias] = sc
		}
	}
}

func getSubCommand(name string) (subCommand, error) {
	if sc, ok := subCommandRegistry.names[name]; ok {
		return sc, nil
	}
	if sc, ok := subCommandRegistry.aliases[name]; ok {
		return sc, nil
	}
	return nil, errUnknownSubCommand(name)
}

func printSubCommandUsage(out io.Writer, sc subCommand, f *flag.FlagSet) {
	aliases := sc.Aliases()
	sort.Strings(aliases)
	aliasesString := strings.Join(aliases, ", ")
	if aliasesString != "" {
		aliasesString = ", " + aliasesString
	}
	nameDescription := wordWrap(sc.Name()+aliasesString+" - "+strings.TrimSpace(sc.Description()), 120)
	fmt.Fprintf(out, "%v\n\nUsage: %v", nameDescription, sc.Name())

	params, paramDesc, hasOptions := sc.Usage()

	if params != "" {
		fmt.Fprintf(out, " %v", params)
	}
	if hasOptions {
		fmt.Fprintf(out, " %v", formatOptionalArgument(commandOptionsName, true))
	}
	fmt.Fprintf(out, "\n\n")

	if params == "" {
		fmt.Fprintf(out, "%v has no %v\n", sc.Name(), formatOptionalArgument(commandParametersName, true))
	} else {
		fmt.Fprintf(out, "%v: %v\n", commandParametersName, paramDesc)
	}

	if hasOptions {
		fmt.Fprintf(out, "%v:\n", commandOptionsName)
		f.SetOutput(out)
		f.PrintDefaults()
	} else {
		fmt.Fprintf(out, "%v has no %v\n", sc.Name(), formatOptionalArgument(commandOptionsName, true))
	}
}

func printSubCommandsUsage(out io.Writer) {
	names := make([]string, 0, len(subCommandRegistry.names))
	for name := range subCommandRegistry.names {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintf(out, "Available %vs:\n", subCommandName)
	for _, name := range names {
		sc := subCommandRegistry.names[name]
		nameAliases := []string{name}
		nameAliases = append(nameAliases, sc.Aliases()...)
		sort.Strings(nameAliases[1:])
		nameAliasesOut := strings.Join(nameAliases, ", ")

		fmt.Fprintf(out, "  %v%v%v\n", nameAliasesOut, pad(16, nameAliasesOut), sc.Synopsis())
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

func wordWrap(text string, columns int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return ""
	}

	buffer := bytes.NewBuffer([]byte(words[0]))
	charsLeft := columns - buffer.Len()
	for _, word := range words[1:] {
		if len(word)+1 > charsLeft {
			buffer.WriteString("\n" + word)
			charsLeft = columns - len(word)
		} else {
			buffer.WriteString(" " + word)
			charsLeft -= 1 + len(word)
		}
	}
	return buffer.String()
}
