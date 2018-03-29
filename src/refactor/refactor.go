package refactor

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
