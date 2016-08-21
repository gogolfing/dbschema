package refactor

import "testing"

func TestChangeSet_UnmarshalXML(t *testing.T) {
	source := `
		<ChangeSet id="id" name="name" author="author">
			<RawSql>raw</RawSql>
			<!-- comment should be ignored -->
			<! directive should be ignored >
			<?proc inst should be ignored?>
			<RawSql>raw</RawSql>
		</ChangeSet>
	`

	cs := &ChangeSet{}
	decodeSourceIntoValue(t, cs, source)
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
