package vars

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type DoesNotExistError string

func (e DoesNotExistError) Error() string {
	return fmt.Sprintf("dbschema/vars: variable %q does not exist", string(e))
}

type ErrInvalidReference string

func (e ErrInvalidReference) Error() string {
	return fmt.Sprintf("dbschema/vars: invalid reference %q", string(e))
}

var referenceRegexp = regexp.MustCompile(`^\$\{[^${}]+\}$`)

type Variables struct {
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

func (v *Variables) ensureValuesExist() {
	if v.values == nil {
		v.values = map[string]string{}
	}
}

func (v *Variables) Put(name, value string) {
	v.ensureValuesExist()
	v.values[name] = value
}

func (v *Variables) Dereference(expr string) (string, error) {
	origExpr := expr
	expr = strings.TrimSpace(expr)
	if !IsVariableReference(expr) {
		return "", ErrInvalidReference(origExpr)
	}
	name := InnerVariableName(expr)
	value, ok := v.GetOk(name)
	if ok {
		return value, nil
	}
	return DereferenceEnv(expr)
}

func (v *Variables) GetOk(name string) (string, bool) {
	value, ok := v.values[name]
	return value, ok
}

func DereferenceEnv(expr string) (string, error) {
	origExpr := expr
	expr = strings.TrimSpace(expr)
	if !IsVariableReference(expr) {
		return "", ErrInvalidReference(origExpr)
	}
	name := InnerVariableName(expr)
	value, ok := GetEnvOk(name)
	if ok {
		return value, nil
	}
	return "", DoesNotExistError(expr)
}

func GetEnvOk(name string) (string, bool) {
	value := os.Getenv(name)
	return value, value != ""
}

func InnerVariableName(expr string) string {
	if !IsVariableReference(expr) {
		return expr
	}
	return expr[2 : len(expr)-1]
}

func IsVariableReference(expr string) bool {
	return referenceRegexp.MatchString(strings.TrimSpace(expr))
}
