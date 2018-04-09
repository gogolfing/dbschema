package scope

import "github.com/gogolfing/dbschema/src/refactor/scope/expand"

type RefMapper interface {
	MapRef(name string) (value string, exists bool)
}

type RefMapperFunc func(string) (string, bool)

func (m RefMapperFunc) MapRef(name string) (string, bool) {
	return m(name)
}

type Scope struct {
	parent *Scope

	refMapper RefMapper
}

func Empty() *Scope {
	return &Scope{
		parent:    nil,
		refMapper: &mapRefMapper{m: nil},
	}
}

func (s *Scope) Parent() *Scope {
	if s.parent == nil {
		return Empty()
	}
	return s.parent
}

func (s *Scope) WithChildRefMapper(m RefMapper) *Scope {
	return &Scope{
		parent:    s,
		refMapper: m,
	}
}

func (s *Scope) Expand(value string) (result string, err error) {
	refMapper := s.refMapper
	if refMapper == nil {
		refMapper = &mapRefMapper{m: nil}
	}

	result, err = expand.Expand(value, refMapper.MapRef)
	if err == nil {
		return
	}
	if expand.IsReferenceDoesNotExistError(err) && s.parent != nil {
		return s.parent.Expand(value)
	}
	return
}

type mapRefMapper struct {
	m map[string]string
}

func (m *mapRefMapper) MapRef(value string) (string, bool) {
	v, ok := m.m[value]
	return v, ok
}
