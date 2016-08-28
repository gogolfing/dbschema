package cli

import (
	"fmt"
	"strings"
)

type connParams [][2]string

func newConnParams() connParams {
	return connParams([][2]string{})
}

func (c *connParams) String() string {
	return fmt.Sprint([][2]string(*c))
}

func (c *connParams) Set(in string) error {
	indexEqual := strings.Index(in, "=")
	if indexEqual < 0 {
		return fmt.Errorf("%v does not contain a key, value pair", in)
	}
	key, value := in[:indexEqual], in[indexEqual+1:]
	*c = append(*c, [2]string{key, value})
	return nil
}

func (c *connParams) eachParamValue(visitor func(name, value string)) {
	for _, pair := range *c {
		visitor(pair[0], pair[1])
	}
}
