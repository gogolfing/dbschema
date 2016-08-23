package refactor

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestChangeSet_UnmarshalXML(t *testing.T) {
	source := `
		<ChangeSet id="id" name="name" author="author">
			<RawSql></RawSql>
			<!-- comment should be ignored -->
			<! directive should be ignored >
			<?proc inst should be ignored?>
			<RawSql></RawSql>
		</ChangeSet>
	`

	cs := &ChangeSet{}
	decodeSourceIntoValue(t, cs, source)

	if !reflect.DeepEqual(cs.XMLName, xml.Name{"", "ChangeSet"}) {
		t.Fail()
	}
	if cs.Id != "id" || *cs.Name != "name" || *cs.Author != "author" {
		t.Fail()
	}
	if !reflect.DeepEqual(cs.changers, []Changer{
		&RawSql{XMLName: xml.Name{Local: "RawSql"}},
		&RawSql{XMLName: xml.Name{Local: "RawSql"}},
	}) {
		t.Fail()
	}
}

func TestChangeSet_UnmarshalXML_errorUnknownTokenType(t *testing.T) {
	source := `
		<ChangeSet id="id" name="name" author="author">
			<RawSql>raw</RawSql>
			<? proc inst invalid ?>
		</ChangeSet>
	`

	cs := &ChangeSet{}
	err := decodeSourceIntoValueError(t, cs, source)

	if err != errUnknownTokenType {
		t.Fail()
	}
}

func TestChangeSet_UnmarshalXML_errorInvalidChanger(t *testing.T) {
}
