package refactor

import (
	"encoding/xml"
	"io"
	"os"
)

type ChangeLog struct {
	XMLName xml.Name `xml:"ChangeLog"`
}

func NewChangeLogFile(path string) (*ChangeLog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := NewChangeLog(file)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return c, nil
}

func NewChangeLog(in io.Reader) (*ChangeLog, error) {
	dec := xml.NewDecoder(in)
	c := &ChangeLog{}
	err := dec.Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ChangeLog) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	return nil
}
