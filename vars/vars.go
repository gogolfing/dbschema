package vars

import (
	"encoding/xml"
	"errors"
	"os"
	"regexp"
)

var ErrDoesNotExist = errors.New("vars: does not exist")

var ErrInvalidReference = errors.New("vars: invalid reference")

var envReferenceRegexp = regexp.MustCompile(`^\$\{[^{}]+\}$`)
var referenceRegexp = regexp.MustCompile(`^\{[^{}]+\}$`)

type Variables struct {
	XMLName xml.Name `xml:"Variables"`

	Values []*Variable `xml:"Variable"`
}

func (v *Variables) Dereference(in string) (string, error) {
	if IsVariableReference(in) {
		name := in[1 : len(in)-1]
		return v.Get(name)
	}
	return DereferenceEnv(in)
}

func DereferenceEnv(in string) (string, error) {
	if !IsEnvVariableReference(in) {
		return "", ErrInvalidReference
	}
	name := in[2 : len(in)-1]
	value := os.Getenv(name)
	if value == "" {
		return "", ErrDoesNotExist
	}
	return value, nil
}

func (v *Variables) Get(name string) (string, error) {
	for _, value := range v.Values {
		if value.Name == name {
			return value.Value, nil
		}
	}
	return "", ErrDoesNotExist
}

type Variable struct {
	XMLName xml.Name `xml:"Variable"`

	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func IsEnvVariableReference(in string) bool {
	return envReferenceRegexp.MatchString(in)
}

func IsVariableReference(in string) bool {
	return referenceRegexp.MatchString(in)
}
