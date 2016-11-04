package refactor

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/gogolfing/dbschema/dialect"
)

type AddPrimaryKey struct {
	XMLName xml.Name `xml:"AddPrimaryKey"`

	Table *StringAttr `xml:"table,attr"`
	Name  *StringAttr `xml:"name,attr"`

	IndexColumns []*IndexColumn `xml:"IndexColumn"`
}

func (apk *AddPrimaryKey) Validate() error {
	return ValidateAll(
		apk.Table.NotEmptyValidator("AddPrimaryKey.table"),
		IndexColumnsValidatorPrimaryKey(apk.IndexColumns),
	)
}

func (apk *AddPrimaryKey) Stmts(ctx Context) ([]*Stmt, error) {
	return StmtsFunc(apk.stmts).Validated(apk, ctx)
}

func (apk *AddPrimaryKey) stmts(ctx Context) ([]*Stmt, error) {
	expanded, err := ExpandAll(
		ctx,
		apk.Table,
		apk.Name,
	)
	if err != nil {
		return nil, err
	}
	table, name := expanded[0], apk.getName(expanded[0], expanded[1])

	columnNamesClause, err := apk.columnNamesClause(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return []*Stmt{
		NewStmtFmt("%v %v %v %v %v %v (%v)",
			dialect.AlterTable,
			ctx.QuoteRef(table),
			dialect.Add,
			dialect.Constraint,
			ctx.QuoteRef(name),
			dialect.PrimaryKey,
			columnNamesClause,
		),
	}, nil
}

func (apk *AddPrimaryKey) getName(table, name string) string {
	if name == "" {
		return fmt.Sprintf("%v_pkey", table)
	}
	return name
}

func (apk *AddPrimaryKey) columnNamesClause(ctx Context) (string, error) {
	names, err := apk.columnNames(ctx)
	if err != nil {
		return "", err
	}
	refNames := []string{}
	for _, name := range names {
		refNames = append(refNames, ctx.QuoteRef(name))
	}
	return strings.Join(refNames, ", "), nil
}

func (apk *AddPrimaryKey) columnNames(ctx Context) ([]string, error) {
	result := []string{}
	for _, ic := range apk.IndexColumns {
		name, err := ic.Name.Expand(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	return result, nil
}
