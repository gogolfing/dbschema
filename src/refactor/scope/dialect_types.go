package scope

import (
	"reflect"
	"strings"

	"github.com/gogolfing/dbschema/src/dialect"
)

const (
	DialectTypesReferencePrefix = "Dialect.Types."
)

type dialectTypesRefMapper struct {
	dialect.Types
}

func NewRefMapperDialectTypes(types dialect.Types) RefMapper {
	return &dialectTypesRefMapper{
		Types: types,
	}
}

func (m *dialectTypesRefMapper) MapRef(name string) (result string, exists bool) {
	if !strings.HasPrefix(name, DialectTypesReferencePrefix) {
		return
	}
	methodName := name[len(DialectTypesReferencePrefix):]

	v := reflect.ValueOf(m.Types)
	method := v.MethodByName(methodName)
	if method.Kind() == reflect.Invalid {
		return
	}

	defer func() {
		recover()
	}()

	out := method.Call([]reflect.Value{})
	result, exists = out[0].Interface().(string), true
	return
}
