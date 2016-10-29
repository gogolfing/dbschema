package dialect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gogolfing/dbschema/conn"
	"github.com/gogolfing/dbschema/vars"
)

//SQL constructs that are the same accross supported dialects.
const (
	CreateTable = "CREATE TABLE"
	AlterTable  = "ALTER TABLE"
	Add         = "ADD"
	Constraint  = "CONSTRAINT"
	PrimaryKey  = "PRIMARY KEY"
	NotNull     = "NOT NULL"
	Default     = "DEFAULT"
	IfNotExists = "IF NOT EXISTS"
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
	ConnectionString(conn *conn.Connection) (string, error)

	DBMS() string

	QuoteRef(in string) string
	QuoteConst(in string) string

	EscapeConst(in string) (string, bool)

	Cast(in, t string) string

	//Following are the "variable" methods.

	Integer() (string, error)
	Int8() (string, error)
	Int16() (string, error)
	Int32() (string, error)
	Int64() (string, error)

	Float32() (string, error)
	Float64() (string, error)

	Char32() (string, error)
	Char64() (string, error)
	Char128() (string, error)
	Char256() (string, error)

	VarChar32() (string, error)
	VarChar64() (string, error)
	VarChar128() (string, error)
	VarChar256() (string, error)
	VarChar512() (string, error)
	VarChar1024() (string, error)

	Text() (string, error)

	Timestamp() (string, error)
	TimestampTz() (string, error)
	Time() (string, error)
	TimeTz() (string, error)
	Date() (string, error)

	Bool() (string, error)

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

const DialectVariablePrefix = "Dialect."

func Expand(expr string, v *vars.Variables, d Dialect) (string, error) {
	origExpr := expr

	if vars.IsVariableReference(expr) {
		value, err := v.Dereference(expr)
		//err will not be vars.ErrInvalidReference because of the prior vars.IsVariableReference() check.
		if err == nil {
			return value, nil
		}

		name := vars.InnerVariableName(expr)
		if !strings.HasPrefix(name, DialectVariablePrefix) {
			//not in v and not a Dialect variable.
			return "", vars.ErrDoesNotExist(origExpr)
		} else {
			name = strings.TrimPrefix(name, DialectVariablePrefix)
		}

		fmt.Println("inner name", name)

		value, err = CallVariableMethodOnDialect(d, name)
		if err != nil {
			return "", vars.ErrDoesNotExist(origExpr)
		}
		return value, nil
	}
	return expr, nil
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

func NewDefaultEscapes() map[string]string {
	return map[string]string{
		"\b": `\b`,
		"\f": `\f`,
		"\n": `\n`,
		"\r": `\r`,
		"\t": `\t`,
	}
}

func DoubleColonCaster(in, t string) string {
	return fmt.Sprintf("%v%v%v", in, "::", t)
}

type DialectStruct struct {
	DBMSValue string

	QuoteRefValue   string
	QuoteConstValue string

	Escapes map[string]string

	Caster func(value, t string) string

	IntegerValue string
	Int8Value    string
	Int16Value   string
	Int32Value   string
	Int64Value   string

	Float32Value string
	Float64Value string

	Char32Value  string
	Char64Value  string
	Char128Value string
	Char256Value string

	VarChar32Value   string
	VarChar64Value   string
	VarChar128Value  string
	VarChar256Value  string
	VarChar512Value  string
	VarChar1024Value string

	TextValue string

	TimestampValue   string
	TimestampTzValue string
	TimeValue        string
	TimeTzValue      string
	DateValue        string

	BoolValue string

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

func (d *DialectStruct) EscapeConst(in string) (string, bool) {
	contains := false
	for key, value := range d.Escapes {
		contains = strings.Contains(in, key) || contains
		in = strings.Replace(in, key, value, -1)
	}
	return d.QuoteConst(in), contains
}

func (d *DialectStruct) Cast(in, t string) string {
	return d.Caster(in, t)
}

func (d *DialectStruct) Integer() (string, error) {
	return d.validString(d.IntegerValue)
}

func (d *DialectStruct) Int8() (string, error) {
	return d.validString(d.Int8Value)
}

func (d *DialectStruct) Int16() (string, error) {
	return d.validString(d.Int16Value)
}

func (d *DialectStruct) Int32() (string, error) {
	return d.validString(d.Int32Value)
}

func (d *DialectStruct) Int64() (string, error) {
	return d.validString(d.Int64Value)
}

func (d *DialectStruct) Float32() (string, error) {
	return d.validString(d.Float32Value)
}

func (d *DialectStruct) Float64() (string, error) {
	return d.validString(d.Float64Value)
}

func (d *DialectStruct) Char32() (string, error) {
	return d.validString(d.Char32Value)
}

func (d *DialectStruct) Char64() (string, error) {
	return d.validString(d.Char64Value)
}

func (d *DialectStruct) Char128() (string, error) {
	return d.validString(d.Char128Value)
}

func (d *DialectStruct) Char256() (string, error) {
	return d.validString(d.Char256Value)
}

func (d *DialectStruct) VarChar32() (string, error) {
	return d.validString(d.VarChar32Value)
}

func (d *DialectStruct) VarChar64() (string, error) {
	return d.validString(d.VarChar64Value)
}

func (d *DialectStruct) VarChar128() (string, error) {
	return d.validString(d.VarChar128Value)
}

func (d *DialectStruct) VarChar256() (string, error) {
	return d.validString(d.VarChar256Value)
}

func (d *DialectStruct) VarChar512() (string, error) {
	return d.validString(d.VarChar512Value)
}

func (d *DialectStruct) VarChar1024() (string, error) {
	return d.validString(d.VarChar1024Value)
}

func (d *DialectStruct) Text() (string, error) {
	return d.validString(d.TextValue)
}

func (d *DialectStruct) Timestamp() (string, error) {
	return d.validString(d.TimestampValue)
}

func (d *DialectStruct) TimestampTz() (string, error) {
	return d.validString(d.TimestampTzValue)
}

func (d *DialectStruct) Time() (string, error) {
	return d.validString(d.TimestampValue)
}

func (d *DialectStruct) TimeTz() (string, error) {
	return d.validString(d.TimeTzValue)
}

func (d *DialectStruct) Date() (string, error) {
	return d.validString(d.DateValue)
}

func (d *DialectStruct) Bool() (string, error) {
	return d.validString(d.BoolValue)
}

//UUID returns d.UUIDValue and ErrNotSupported if it is empty.
func (d *DialectStruct) UUID() (string, error) {
	return d.validString(d.UUIDValue)
}

func (d *DialectStruct) validString(value string) (string, error) {
	if value == "" {
		return "", ErrNotSupported
	}
	return value, nil
}

func Quote(in, q string) string {
	return fmt.Sprintf("%v%v%v", q, in, q)
}
