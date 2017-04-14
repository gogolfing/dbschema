package refactor

type Expander interface {
	Expand(Context) (string, error)
}

func ExpandAll(ctx Context, expanders ...Expander) ([]string, error) {
	result := []string{}
	for _, expander := range expanders {
		value, err := expander.Expand(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

type ExpanderFunc func(Context) (string, error)

func (ef ExpanderFunc) Expand(ctx Context) (string, error) {
	return ef(ctx)
}
