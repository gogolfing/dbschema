package dialect

type Types interface {
	Integer() string
	Int8() string
	Int16() string
	Int32() string
	Int64() string

	Float32() string
	Float64() string

	Char32() string
	Char64() string
	Char128() string
	Char256() string

	VarChar32() string
	VarChar64() string
	VarChar128() string
	VarChar256() string
	VarChar512() string
	VarChar1024() string

	Text() string

	ByteArray() string

	Timestamp() string
	TimestampTz() string
	Time() string
	TimeTz() string
	Date() string

	Bool() string

	UUID() string
}

type TypesStruct struct {
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

	ByteArrayValue string

	TimestampValue   string
	TimestampTzValue string
	TimeValue        string
	TimeTzValue      string
	DateValue        string

	BoolValue string

	UUIDValue string
}

func (t *TypesStruct) Integer() string {
	return t.IntegerValue
}

func (t *TypesStruct) Int8() string {
	return t.Int8Value
}

func (t *TypesStruct) Int16() string {
	return t.Int16Value
}

func (t *TypesStruct) Int32() string {
	return t.Int32Value
}

func (t *TypesStruct) Int64() string {
	return t.Int64Value
}

func (t *TypesStruct) Float32() string {
	return t.Float32Value
}

func (t *TypesStruct) Float64() string {
	return t.Float64Value
}

func (t *TypesStruct) Char32() string {
	return t.Char32Value
}

func (t *TypesStruct) Char64() string {
	return t.Char64Value
}

func (t *TypesStruct) Char128() string {
	return t.Char128Value
}

func (t *TypesStruct) Char256() string {
	return t.Char256Value
}

func (t *TypesStruct) VarChar32() string {
	return t.VarChar32Value
}

func (t *TypesStruct) VarChar64() string {
	return t.VarChar64Value
}

func (t *TypesStruct) VarChar128() string {
	return t.VarChar128Value
}

func (t *TypesStruct) VarChar256() string {
	return t.VarChar256Value
}

func (t *TypesStruct) VarChar512() string {
	return t.VarChar512Value
}

func (t *TypesStruct) VarChar1024() string {
	return t.VarChar1024Value
}

func (t *TypesStruct) Text() string {
	return t.TextValue
}

func (t *TypesStruct) ByteArray() string {
	return t.ByteArrayValue
}

func (t *TypesStruct) Timestamp() string {
	return t.TimestampValue
}

func (t *TypesStruct) TimestampTz() string {
	return t.TimestampTzValue
}

func (t *TypesStruct) Time() string {
	return t.TimeValue
}

func (t *TypesStruct) TimeTz() string {
	return t.TimeTzValue
}

func (t *TypesStruct) Date() string {
	return t.DateValue
}

func (t *TypesStruct) Bool() string {
	return t.BoolValue
}

func (t *TypesStruct) UUID() string {
	return t.UUIDValue
}
