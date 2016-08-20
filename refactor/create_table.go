package refactor

// import (
// 	"encoding/xml"
// 	"fmt"
// )

// type CreateTable struct {
// 	XMLName xml.Name `xml:"CreateTable"`

// 	Name string `xml:"name,attr"`

// 	Columns []*Column `xml:"Columns"`
// }

// func (a *CreateTable) Up(ctx *Context) (stmts []string, err error) {
// 	if err := a.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return StmtsFromFuncs(ctx, a.upTable)
// }

// func (a *CreateTable) upTable(ctx *Context) (string, error) {
// 	columns, err := a.upColumns(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprintf("%v %v (%v)", ctx.CreateTable, a.Name, columns), nil
// }

// func (a *CreateTable) upColumns(ctx *Context) (string, error) {
// 	return "", nil
// }

// func (a *CreateTable) Validate() error {
// 	if a.Name == "" {
// 		return ErrInvalid("CreateTable.Name cannot be empty")
// 	}
// 	return nil
// }
