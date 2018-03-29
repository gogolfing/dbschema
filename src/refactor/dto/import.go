package dto

import (
	"encoding/xml"
	"fmt"
)

type InvalidImportPathError string

func (e InvalidImportPathError) Error() string {
	return fmt.Sprintf("dbschema/refactor/dto: invalid import path %q", string(e))
}

type Import struct {
	XMLName xml.Name `xml:"Import"`

	Path string `xml:"path,attr"`
}

func newImport() *Import {
	return &Import{}
}

func (i *Import) Validate() error {
	if i.Path == "" {
		return InvalidImportPathError(i.Path)
	}

	return nil
}
