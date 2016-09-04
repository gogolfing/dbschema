package cli

import "github.com/gogolfing/dbschema/dbschema"

type subCommand interface {
	Run(*dbschema.DbSchema, dbschema.Driver)
}

func (c *CLI) parseSubCommandFromName(name string, args []string) (subCommand, error) {
	return nil, nil
}
