package refactor

import (
	"encoding/xml"
	"io"

	"github.com/gogolfing/dbschema/refactor/dto"
	"github.com/gogolfing/dbschema/vars"
)

const (
	errInvalidImportPath = ErrInvalid("Import.path cannot be empty")
)

type ChangeLog struct {
	XMLName xml.Name `xml:"ChangeLog"`

	ChangeLogTableName     *StringAttr `xml:"changeLogTableName,attr"`
	ChangeLogLockTableName *StringAttr `xml:"changeLogLockTableName,attr"`

	path string

	Variables *vars.Variables `xml:"Variables"`

	ChangeSets []*ChangeSet `xml:"ChangeSet"`
}

func NewChangeLogFile(path string) (*ChangeLog, error) {
	return dto.NewChangeLogFile(path)
}

func NewChangeLogReader(path string, in io.Reader) (*ChangeLog, error) {
	dec := xml.NewDecoder(in)
	c := &ChangeLog{}
	c.path = path
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ChangeLog) ChangeSetsSubSlice(afterId, beforeId string) []*ChangeSet {
	afterIndex := 0
	if afterId != "" {
		afterIndex = c.findChangeSetIdIndex(afterIndex, afterId)
	}
	beforeIndex := c.findChangeSetIdIndex(afterIndex+1, beforeId)

	result := []*ChangeSet{}
	for index := afterIndex + 1; index < len(c.ChangeSets) && index < beforeIndex; index++ {
		result = append(result, c.ChangeSets[index])
	}

	return result
}

func (c *ChangeLog) findChangeSetIdIndex(startIndex int, id string) int {
	for startIndex < len(c.ChangeSets) && c.ChangeSets[startIndex].Id != id {
		startIndex++
	}
	return startIndex
}
