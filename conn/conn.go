package conn

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/gogolfing/dbschema/vars"
)

type Connection struct {
	XMLName xml.Name `xml:"Connection"`

	DBMS string `xml:"dbms,attr"`

	Host string `xml:"host,attr"`

	Port int `xml:"port,attr"`

	User string `xml:"user,attr"`

	Password string `xml:"password,attr"`

	params map[string]string
}

func NewConnectionFile(path string) (*Connection, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewConnectionReader(file)
}

func NewConnectionReader(in io.Reader) (*Connection, error) {
	dec := xml.NewDecoder(in)
	c := &Connection{}
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}

/*
func (c *Connection) DBMSValue() (string, error) {
	return possibleVariable(c.DBMS)
}

func (c *Connection) HostValue() (string, error) {
	return possibleVariable(c.Host)
}

func (c *Connection) UserValue() (string, error) {
	return possibleVariable(c.User)
}

func (c *Connection) PasswordValue() (string, error) {
	return possibleVariable(c.Password)
}
*/

func (c *Connection) PutParam(name, value string) {
	c.ensureParamsExist()
	c.params[name] = value
}

func (c *Connection) EachParamValue(visitor func(name, value string)) {
	for name, value := range c.params {
		visitor(name, value)
	}
}

func (c *Connection) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	c.ensureParamsExist()
	xmlC := &xmlConnection{}
	err := dec.DecodeElement(xmlC, &start)
	if err != nil {
		return err
	}

	c.XMLName = xmlC.XMLName

	c.DBMS, err = possibleVariable(xmlC.DBMS)
	if err != nil {
		return err
	}
	c.Host = xmlC.Host
	c.Port = xmlC.Port
	c.User = xmlC.User
	c.Password = xmlC.Password

	for _, param := range xmlC.Params {
		c.PutParam(param.Name, param.Value)
	}
	return nil
}

func (c *Connection) ensureParamsExist() {
	if c.params == nil {
		c.params = map[string]string{}
	}
}

type xmlConnection struct {
	XMLName xml.Name `xml:"Connection"`

	DBMS string `xml:"dbms,attr"`

	Host string `xml:"host,attr"`

	Port int `xml:"port,attr"`

	User string `xml:"user,attr"`

	Password string `xml:"password,attr"`

	Params []*Param `xml:"Param"`
}

type Param struct {
	XMLName xml.Name `xml:"Param"`

	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func possibleVariable(expr string) (string, error) {
	if vars.IsEnvVariableReference(expr) {
		return vars.DereferenceEnv(expr)
	}
	return expr, nil
}
