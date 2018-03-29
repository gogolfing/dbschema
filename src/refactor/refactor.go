package refactor

import (
	"github.com/gogolfing/dbschema/src/refactor/dto"
	"github.com/gogolfing/dbschema/src/vars"
)

type NullString struct {
	String string
	Valid  bool
}

func NewNullString(s *string) NullString {
	ns := NullString{
		Valid: s != nil,
	}
	if ns.Valid {
		ns.String = *s
	}
	return ns
}

func (ns *NullString) DefaultExpander(def string) Expander {
	return ExpanderFunc(func(ctx Context) (string, error) {
		if !ns.Valid || len(ns.String) == 0 {
			return def, nil
		}
		return ctx.Expand(ns.String)
	})
}

type NullBool struct {
	Raw   string
	Bool  bool
	Valid bool
}

func (nb *NullBool) DefaultExpander(def bool) Expander {
	return ExpanderFunc(func(ctx Context) (string, error) {
		return "", nil
	})
}

func newChangeLogDto(dtoCl *dto.ChangeLog) (*ChangeLog, error) {
	changeSets, err := newChangeSetsDto(dtoCl.ChangeSets)
	if err != nil {
		return nil, err
	}

	cl := &ChangeLog{
		TableName:     NewNullString(dtoCl.TableName),
		LockTableName: NewNullString(dtoCl.LockTableName),
		Variables:     newVariablesDto(dtoCl.Variables),
		ChangeSets:    changeSets,
	}
	return cl, nil
}

func newVariablesDto(dtoVars *dto.Variables) *vars.Variables {
	v := &vars.Variables{}
	for _, dtoVar := range dtoVars.Values {
		v.Put(dtoVar.Name, dtoVar.Value)
	}
	return v
}

func newChangeSetsDto(dtoCss []*dto.ChangeSet) ([]*ChangeSet, error) {
	result := []*ChangeSet{}
	for _, dtoCs := range dtoCss {
		cs, err := newChangeSetDto(dtoCs)
		if err != nil {
			return nil, err
		}
		result = append(result, cs)
	}
	return result, nil
}

func newChangeSetDto(dtoCs *dto.ChangeSet) (*ChangeSet, error) {
	changers, err := newChangersDto(dtoCs.Changers)
	if err != nil {
		return nil, err
	}

	return &ChangeSet{
		Id:       dtoCs.Id,
		Name:     NewNullString(dtoCs.Name),
		Author:   NewNullString(dtoCs.Author),
		Changers: changers,
	}, nil
}

func newChangersDto(dtoChangers []dto.Changer) ([]Changer, error) {
	result := []Changer{}
	for _, dtoChanger := range dtoChangers {
		changer, err := newChangerDto(dtoChanger)
		if err != nil {
			return nil, err
		}
		result = append(result, changer)
	}
	return result, nil
}

func newChangerDto(dtoChanger dto.Changer) (Changer, error) {
	switch dtoType := dtoChanger.(type) {
	case *dto.RawSql:
		return newRawSqlDto(dtoType), nil
	}
	return nil, &UnknownDTOChangerError{dtoChanger}
}

func newRawSqlDto(rawSql *dto.RawSql) Changer {
	return &RawSql{
		UpStmts:   newStmtsDto(rawSql.Up.Stmts),
		DownStmts: newStmtsDto(rawSql.Down.Stmts),
		dto:       rawSql,
	}
}

func newStmtsDto(dtoStmts []*dto.Stmt) []*Stmt {
	result := make([]*Stmt, 0, len(dtoStmts))
	for _, stmt := range dtoStmts {
		result = append(result, &Stmt{Raw: stmt.Raw})
	}
	return result
}
