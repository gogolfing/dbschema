package dialect

import (
	"errors"
	"fmt"
	"reflect"
)

const TagVarname = "varname"

var ErrFieldDoesNotExist = errors.New("dialect: field does not exist")

type Dialect struct {
	CreateTable string

	QuoteRef   string
	QuoteConst string

	NotNull string
	Default string

	Int  string `varname:"Int"`
	UUID string `varname:"UUID"`
}

func (d *Dialect) ValueOfVariableField(varname string) (string, error) {
	value := reflect.ValueOf(*d)
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		if fieldVarname := field.Tag.Get(TagVarname); fieldVarname == varname {
			return fmt.Sprint(value.Field(i).Interface()), nil
		}
	}
	return "", ErrFieldDoesNotExist
}
