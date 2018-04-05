package refactor

type ChangeSet struct {
	Id string

	Name   NullString
	Author NullString

	Tags []string

	Changers []Changer
}

/*
func (c *ChangeSet) Sha256Sum() ([]byte, error) {
	hash := sha256.New()
	enc := xml.NewEncoder(hash)
	for _, changer := range c.Changers {
		if err := dto.EncodeRefactorChangerXML(enc, changer); err != nil {
			return err
		}
	}
	return hash.Sum(nil), nil
}
*/
