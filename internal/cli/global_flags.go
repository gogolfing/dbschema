package cli

import (
	"flag"
	"fmt"
	"os/user"
	"strings"
)

type globalFlags struct {
	verbose       bool
	connPath      string
	changeLogPath string
	dbms          string
	host          string
	port          int
	user          string
	password      string
	database      string
	connParams    connParams
}

func (gf *globalFlags) SetFlags(f *flag.FlagSet) {
	u, err := user.Current()
	if err != nil {
		u = &user.User{Username: ""}
	}

	f.BoolVar(&gf.verbose, "v", false, "print verbose output")
	f.StringVar(&gf.connPath, "conn", "connection.xml", "path to connection file")
	f.StringVar(&gf.changeLogPath, "changelog", "changelog.xml", "path to change log file")
	f.StringVar(&gf.dbms, "dbms", "", "the type of the dbms to connect to. this will override the value in -conn if not empty")
	f.StringVar(&gf.host, "host", "", "host to connect to. this will override the value in -conn if not empty")
	f.IntVar(&gf.port, "port", 0, "port to connect to. this will override the value in -conn if not zero")
	f.StringVar(&gf.user, "user", u.Username, "the user to connect to as. this will override the value in -conn if not empty")
	f.StringVar(&gf.password, "password", "", "password to connect with. this will override the value in -conn if not empty")
	f.StringVar(&gf.database, "database", "", "database to connect to. this will override the value in -conn if not empty")
	f.Var(&gf.connParams, "conn-param", "list of connection parameters in the form <name>=<value>. should be set with multiple flag definitions. these will override already set parameters in -conn")
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
