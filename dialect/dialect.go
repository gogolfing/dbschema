package dialect

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrFieldDoesNotExist = errors.New("dialect: field does not exist")

type Dialect struct {
	//Actions.

	CreateTable string

	//Types.

	Int  string
	UUID string
}

func (d *Dialect) ValueOfField(name string) (string, error) {
	value := reflect.ValueOf(*d)
	fieldValue := value.FieldByName(name)
	if fieldValue.Kind() == reflect.Invalid {
		return "", ErrFieldDoesNotExist
	}
	return fmt.Sprint(fieldValue.Interface()), nil
}
