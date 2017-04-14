package refactor

type Context interface {
	Expand(string) (string, error)
}
