package cli

type CreateConnectionError struct {
	Err error
}

func (e *CreateConnectionError) Error() string {
	return e.Err.Error()
}

type CreateDialectError struct {
	Err error
}

func (e *CreateDialectError) Error() string {
	return e.Err.Error()
}

type CreateChangeLogError struct {
	Err error
}

func (e *CreateChangeLogError) Error() string {
	return e.Err.Error()
}

type CreateDBSchemaError struct {
	Err error
}

func (e *CreateDBSchemaError) Error() string {
	return e.Err.Error()
}
