package dto

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

func TestChangeSet_UnmarshalXML_Error_UnknownChangerType(t *testing.T) {
	raw := `
		<ChangeSet>
			<UnknownChangerType> something </UnknownChangerType>
		</ChangeSet>
	`

	dec := xml.NewDecoder(strings.NewReader(raw))

	cs := &ChangeSet{}
	err := dec.Decode(cs)
	if !reflect.DeepEqual(err, UnknownChangerTypeError("UnknownChangerType")) {
		t.Fatal()
	}
}

func TestChangeSet_UnmarshalXML_Ok(t *testing.T) {
	raw := `
		<ChangeSet id="foo" name="bar">
			<RawSql></RawSql>
		</ChangeSet>
	`

	dec := xml.NewDecoder(strings.NewReader(raw))

	cs := &ChangeSet{}
	if err := dec.Decode(cs); err != nil {
		t.Fatal(err)
	}

	if len(cs.Changers) != 1 {
		t.Fatal()
	}
	if _, ok := cs.Changers[0].(*RawSql); !ok {
		t.Fatal()
	}
	if cs.Id != "foo" {
		t.Fatal()
	}
	if *cs.Name != "bar" {
		t.Fatal()
	}
	if cs.Author != nil {
		t.Fatal()
	}
}
