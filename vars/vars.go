package vars

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ErrDoesNotExist string

func (e ErrDoesNotExist) Error() string {
	return fmt.Sprintf("dbschema/vars: variable %q does not exist", string(e))
}

type ErrInvalidReference string

func (e ErrInvalidReference) Error() string {
	return fmt.Sprintf("dbschema/vars: invalid reference %q", string(e))
}

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
	expr = strings.TrimSpace(expr)
	if IsVariableReference(expr) {
		name := expr[1 : len(expr)-1]
		value := v.get(name)
		if value == nil {
			return "", ErrDoesNotExist(expr)
		}
		return *value, nil
	}
	return DereferenceEnv(expr)
}

func DereferenceEnv(expr string) (string, error) {
	expr = strings.TrimSpace(expr)
	if !IsEnvVariableReference(expr) {
		return "", ErrInvalidReference(expr)
	}
	name := expr[2 : len(expr)-1]
	value := os.Getenv(name)
	if value == "" {
		return "", ErrDoesNotExist(expr)
	}
	return value, nil
}

func (v *Variables) get(name string) *string {
	value, ok := v.values[name]
	if !ok {
		return nil
	}
	return &value
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
	return envReferenceRegexp.MatchString(strings.TrimSpace(expr))
}

func IsVariableReference(expr string) bool {
	return referenceRegexp.MatchString(strings.TrimSpace(expr))
}
