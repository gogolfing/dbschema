package refactor

import "encoding/xml"

const (
	errUpMustBeDefined   = ErrInvalid("RawSql > Up must be defined")
	errDownMustDeDefined = ErrInvalid("RawSql > Down must be defined")
)

type RawSql struct {
	XMLName xml.Name `xml:"RawSql"`

	UpValue   *string `xml:"Up"`
	DownValue *string `xml:"Down"`
}

func (r *RawSql) Validate(ctx Context) error {
	if r.UpValue == nil {
		return errUpMustBeDefined
	}
	if r.DownValue == nil {
		return errDownMustDeDefined
	}
	return nil
}

func (r *RawSql) Up(ctx Context) ([]string, error) {
	return []string{*r.UpValue}, nil
}

func (r *RawSql) Down(ctx Context) ([]string, error) {
	return []string{*r.DownValue}, nil
}
