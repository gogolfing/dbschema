package refactor

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNewChangeLogFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "dbschema_ChangeLog")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	importedChangeSetSource := `
		<ChangeSet id="id2" author="author imported" name="name2">
			<RawSql>raw sql imported</RawSql>
		</ChangeSet>
	`
	importedChangeSetFile, err := ioutil.TempFile(dir, "ChangeSet")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, importedChangeSetFile, importedChangeSetSource)

	changeLogFileSource := `
		<ChangeLog>

			<Variables>
				<Variable name="varNameOne" value="varValueOne" />
			</Variables>

			<ChangeSet id="id1" author="author local" name="name1">
				<RawSql>raw sql 1</RawSql>
				<RawSql>raw sql 2</RawSql>
			</ChangeSet>

			<Variables>
				<Variable name="varNameTwo" value="varValueTwo" />
			</Variables>

			<Import path="` + path.Base(importedChangeSetFile.Name()) + `" />

			<ChangeSet id="id3" author="author local" name="name1">
				<RawSql>raw sql 3</RawSql>
			</ChangeSet>

		</ChangeLog>
	`

	changeLogFile, err := ioutil.TempFile(dir, "ChangeLog")
	writeFile(t, changeLogFile, changeLogFileSource)

	changeLog, err := NewChangeLogFile(changeLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	//the following checks are done in separation because we cannot create a
	//Variables the way we need to for reflect.DeepEqual().

	//make sure path is set correctly.
	if changeLog.path != changeLogFile.Name() {
		t.Errorf("changeLog.path = %v WANT %v", changeLog.path, changeLogFile.Name())
	}

	//make sure variables are all what they should be.
	testVariablesEqual(t, changeLog,
		"varNameOne", "varValueOne",
		"varNameTwo", "varValueTwo",
	)

	//make sure ChangeSets are what they should be.
	if len(changeLog.ChangeSets) != 3 {
		t.Fail()
	}
	for i, id := range []string{"id1", "id2", "id3"} {
		if csid := changeLog.ChangeSets[i].Id; csid != id {
			t.Errorf("changeLog.ChangeSets[%v].Id = %v WANT %v", i, csid, id)
		}
	}
}

func writeFile(t *testing.T, file *os.File, source string) {
	if _, err := file.WriteString(source); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func testVariablesEqual(t *testing.T, c *ChangeLog, vars ...string) {
	for i := 0; i < len(vars); i += 2 {
		name, value := vars[i], vars[i+1]
		if actual, err := c.Variables.Get(name); actual != value || err != nil {
			t.Errorf("c.Variables.Get(%v) = %v, %v WANT %v, %v", name, actual, err, value, nil)
		}
	}
	if count := len(vars) >> 1; count != c.Variables.Len() {
		t.Errorf("c.Variables.Len() = %v WANT %v", c.Variables.Len(), count)
	}
}
