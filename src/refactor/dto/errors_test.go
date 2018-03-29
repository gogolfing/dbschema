package dto

import "testing"

func TestUnknownChangerTypeError_Error(t *testing.T) {
	err := UnknownChangerTypeError("type")

	if err.Error() != `dbschema/refactor/dto: unknown Changer type in ChangeSet "type"` {
		t.Fatal()
	}
}

func TestInvalidImportPathError_Error(t *testing.T) {
	err := InvalidImportPathError("")

	if err.Error() != `dbschema/refactor/dto: invalid import path ""` {
		t.Fatal()
	}
}
