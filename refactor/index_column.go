package refactor

import "encoding/xml"

type IndexColumn struct {
	XMLName xml.Name `xml:"IndexColumn"`

	Name          *StringAttr `xml:"name,attr"`
	RefTable      *StringAttr `xml:"refTable,attr"`
	RefColumnName *StringAttr `xml:"refColumnName,attr"`
}

func (ic *IndexColumn) Validate() error {
	return ValidateAll(
		ic.Name.NotEmptyValidator("IndexColumn.name"),
		ic.RefTable.NotEmptyValidator("IndexColumn.refTable"),
	)
}

func IndexColumnsValidatorPrimaryKey(ics []*IndexColumn) Validator {
	return ValidatorFunc(func() error {
		for _, ic := range ics {
			err := ic.Name.NotEmptyValidator("IndexColumn.name").Validate()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func IndexColumnsValidator(ics []*IndexColumn) Validator {
	return ValidatorFunc(func() error {
		for _, ic := range ics {
			if err := ic.Validate(); err != nil {
				return err
			}
		}
		return nil
	})
}
