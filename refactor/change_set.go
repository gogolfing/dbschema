package refactor

import (
	"crypto/sha256"
	"encoding/xml"
)

type ChangeSet struct {
	Id string

	Name   NullString
	Author NullString

	Tags []string

	Changers []Changer
}

func (c *ChangeSet) Sha256Sum() ([]byte, error) {
	hash := sha256.New()
	for _, changer := range c.Changers {
		dto, err := changer.DTO()
		if err != nil {
			return nil, err
		}
		b, err := xml.Marshal(dto)
		if err != nil {
			return nil, err
		}
		hash.Write(b)
	}
	return hash.Sum(nil), nil
}
