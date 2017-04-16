package cli

import (
	"flag"
	"fmt"
	"strings"
)

type globalFlags struct {
	verbose       bool
	changeLogPath string
	conn          string
	dbms          string
}

func newGlobalFlagsDefault() *globalFlags {
	return &globalFlags{
		changeLogPath: "changelog.xml",
		connPath:      "connection.xml",
	}
}

func (gf *globalFlags) SetFlags(f *flag.FlagSet) {
	fmt.Printf("setting flags %p\n", gf)
	f.BoolVar(&gf.verbose, "v", false, "print verbose output")
	f.StringVar(&gf.conn, "conn", "", "path to connection file")
	f.StringVar(&gf.changeLogPath, "changelog", "changelog.xml", "path to change log file")
	f.StringVar(&gf.dbms, "dbms", "", "the type of the dbms to connect to. this will override the DBSCHEMA_DBMS environment variable if set")
}

type connParams [][2]string

func newConnParams() connParams {
	return connParams([][2]string{})
}

func (c *connParams) String() string {
	if c == nil {
		return "[]"
	}
	return fmt.Sprint([][2]string(*c))
}

func (c *connParams) Set(in string) error {
	indexEqual := strings.Index(in, "=")
	if indexEqual < 0 {
		return fmt.Errorf("%s does not contain a key, value pair", in)
	}
	key, value := in[:indexEqual], in[indexEqual+1:]
	*c = append(*c, [2]string{key, value})
	return nil
}

func (c *connParams) eachKeyValue(visitor func(key, value string)) {
	for _, pair := range *c {
		visitor(pair[0], pair[1])
	}
}
