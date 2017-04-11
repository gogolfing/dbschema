package refactor

import "encoding/xml"

const (
	errChangeSetCannotBeEmpty = ErrInvalid("ChangeSet cannot be empty")
)

type ChangeSet struct {
	XMLName xml.Name `xml:"ChangeSet"`

	Id string `xml:"id,attr"`

	Name *string `xml:"name,attr"`

	Author *string `xml:"author,attr"`

	changers []Changer
}

//add validation for an empty changeset.

func (c *ChangeSet) Stmts(ctx Context) (stmts []*Stmt, err error) {
	for _, changer := range c.changers {
		temp, err := changer.Stmts(ctx)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, temp...)
	}
	return
}
