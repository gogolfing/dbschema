package cli

import "github.com/gogolfing/dbschema/dbschema"

type subCommandFunc func(*dbschema.DBSchema, subCommandFlags) error

type subCommandFlags interface{}

func (c *CLI) parseSubCommandFromName(name string, args []string) (subCommandFunc, subCommandFlags, error) {
	return nil, nil, nil
}
