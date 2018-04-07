package dialect

import (
	"fmt"
	"strings"
)

func Quote(in, q string) string {
	return fmt.Sprintf("%v%v%v", q, in, q)
}

func NewDefaultEscapes() map[string]string {
	return map[string]string{
		"\b": `\b`,
		"\f": `\f`,
		"\n": `\n`,
		"\r": `\r`,
		"\t": `\t`,
	}
}

func DoubleColonCaster(in, t string) string {
	return fmt.Sprintf("%v%v%v", in, "::", t)
}

type Syntax interface {
	QuoteRef(in string) string
	QuoteConst(in string) string

	EscapeConst(in string) (string, bool)

	Cast(in, t string) string

	Placeholder(index int) string
}

type SyntaxStruct struct {
	QuoteRefValue   string
	QuoteConstValue string

	Escapes map[string]string

	Caster func(value, t string) string

	PlaceholderValue func(num int) string
}

func (s *SyntaxStruct) QuoteRef(in string) string {
	return Quote(in, s.QuoteRefValue)
}

func (s *SyntaxStruct) QuoteConst(in string) string {
	return Quote(in, s.QuoteConstValue)
}

func (s *SyntaxStruct) EscapeConst(in string) (string, bool) {
	contains := false
	for key, value := range s.Escapes {
		contains = strings.Contains(in, key) || contains
		in = strings.Replace(in, key, value, -1)
	}
	return s.QuoteConst(in), contains
}

func (s *SyntaxStruct) Cast(in, t string) string {
	return s.Caster(in, t)
}

func (s *SyntaxStruct) Placeholder(num int) string {
	return s.PlaceholderValue(num)
}
