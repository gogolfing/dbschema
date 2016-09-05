package dialect

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gogolfing/dbschema/conn"
)

//SQL statements that are the same accross supported dialects.
const (
	CreateTable = "CREATE TABLE"
)

//SQL modifiers that are the same across supported dialects.
const (
	NotNull = "NOT NULL"
	Default = "DEFAULT"
)

var ErrMethodDoesNotExist = errors.New("dbschema/dialect: method does not exist")

var ErrInvalidVariableMethodType = errors.New("dbschema/dialect: invalid method type")

var ErrNotSupported = errors.New("dbschema/dialect: not supported")

type ErrUnsupportedDBMS string

func (e ErrUnsupportedDBMS) Error() string {
	return fmt.Sprintf("dbschema/dialect: unsupported or undefined DBMS %q", string(e))
}

//Dialect is the contract that all database types must implement for dbschema to
//generate the correct SQL for a given DBMS dialect.
type Dialect interface {
	DBMS() string

	ConnectionString(conn *conn.Connection) (string, error)

	QuoteRef(in string) string

	QuoteConst(in string) string

	//Following are the "variable" methods.

	Int() (string, error)

	UUID() (string, error)
}

func NewDialect(dbms string) (Dialect, error) {
	dialects := map[string]func() Dialect{
		Postgresql: NewDialectPostgresql,
	}
	d, ok := dialects[dbms]
	if !ok {
		return nil, ErrUnsupportedDBMS(dbms)
	}
	return d(), nil
}

type dialect struct {
	connectionString func(conn *conn.Connection) (string, error)
	*DialectStruct
}

func newDialect(connStringFunc func(*conn.Connection) (string, error), ds *DialectStruct) Dialect {
	return &dialect{
		connectionString: connStringFunc,
		DialectStruct:    ds,
	}
}

func (d *dialect) ConnectionString(conn *conn.Connection) (string, error) {
	return d.connectionString(conn)
}

func CallVariableMethodOnDialect(d Dialect, name string) (value string, err error) {
	v := reflect.ValueOf(d)
	method := v.MethodByName(name)
	if method.Kind() == reflect.Invalid {
		return "", ErrMethodDoesNotExist
	}
	if !isMethodIsOfVariableType(method) {
		return "", ErrInvalidVariableMethodType
	}
	out := method.Call([]reflect.Value{})
	value = out[0].Interface().(string)
	err, _ = out[1].Interface().(error)
	return
}

func isMethodIsOfVariableType(method reflect.Value) bool {
	t := method.Type()
	if t.NumIn() != 0 {
		return false
	}
	if t.NumOut() != 2 {
		return false
	}
	if t.Out(0).Kind() != reflect.String {
		return false
	}
	if out := t.Out(1); out.Kind() != reflect.Interface || out.String() != "error" {
		return false
	}
	return true
}

type DialectStruct struct {
	DBMSValue string

	QuoteRefValue string

	QuoteConstValue string

	IntValue string

	UUIDValue string
}

func (d *DialectStruct) DBMS() string {
	return d.DBMSValue
}

func (d *DialectStruct) QuoteRef(in string) string {
	return Quote(in, d.QuoteRefValue)
}

func (d *DialectStruct) QuoteConst(in string) string {
	return Quote(in, d.QuoteConstValue)
}

func (d *DialectStruct) Int() (string, error) {
	return d.validString(d.IntValue)
}

func (d *DialectStruct) UUID() (string, error) {
	return d.validString(d.UUIDValue)
}

func (d *DialectStruct) validString(value string) (string, error) {
	if value == "" {
		return "", ErrNotSupported
	}
	return value, nil
}

func Quote(in, quote string) string {
	return in
}
