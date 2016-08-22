package vars

import (
	"encoding/xml"
	"errors"
	"os"
	"regexp"
)

var ErrDoesNotExist = errors.New("dbschema/vars: does not exist")

var ErrInvalidReference = errors.New("dbschema/vars: invalid reference")

var envReferenceRegexp = regexp.MustCompile(`^\$\{[^{}]+\}$`)
var referenceRegexp = regexp.MustCompile(`^\{[^{}]+\}$`)

type Variables struct {
	XMLName xml.Name `xml:"Variables"`

	values map[string]string
}

func (v *Variables) Len() int {
	return len(v.values)
}

func (v *Variables) Merge(other *Variables) {
	v.ensureValuesExist()
	for name, value := range other.values {
		v.values[name] = value
	}
}

func (v *Variables) Put(variable *Variable) {
	v.ensureValuesExist()
	v.values[variable.Name] = variable.Value
}

func (v *Variables) Dereference(expr string) (string, error) {
	if IsVariableReference(expr) {
		name := expr[1 : len(expr)-1]
		return v.Get(name)
	}
	return DereferenceEnv(expr)
}

func DereferenceEnv(expr string) (string, error) {
	if !IsEnvVariableReference(expr) {
		return "", ErrInvalidReference
	}
	name := expr[2 : len(expr)-1]
	value := os.Getenv(name)
	if value == "" {
		return "", ErrDoesNotExist
	}
	return value, nil
}

func (v *Variables) Get(name string) (string, error) {
	value, ok := v.values[name]
	if !ok {
		return "", ErrDoesNotExist
	}
	return value, nil
}

func (v *Variables) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	v.ensureValuesExist()
	xmlV := &xmlVariables{}
	if err := dec.DecodeElement(xmlV, &start); err != nil {
		return err
	}
	v.XMLName = xmlV.XMLName
	for _, variable := range xmlV.Values {
		v.Put(variable)
	}
	return nil
}

func (v *Variables) ensureValuesExist() {
	if v.values == nil {
		v.values = map[string]string{}
	}
}

type xmlVariables struct {
	XMLName xml.Name `xml:"Variables"`

	Values []*Variable `xml:"Variable"`
}

type Variable struct {
	XMLName xml.Name `xml:"Variable"`

	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func IsEnvVariableReference(expr string) bool {
	return envReferenceRegexp.MatchString(expr)
}

func IsVariableReference(expr string) bool {
	return referenceRegexp.MatchString(expr)
}
