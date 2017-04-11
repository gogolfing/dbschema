package dto

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestChangeLog_UnmarshalXML_Error_InvalidImportPath(t *testing.T) {
	raw := `
		<ChangeLog>
			<ChangeSet>
			</ChangeSet>
			<Import path="" />
		</ChangeLog>
	`

	dec := xml.NewDecoder(strings.NewReader(raw))

	cl := &ChangeLog{}
	err := dec.Decode(cl)
	if !reflect.DeepEqual(err, InvalidImportPathError("")) {
		t.Fatal()
	}
}

func TestChangeLog_UnmarshalXML_Error_Ok(t *testing.T) {
	tempDir, _ := ioutil.TempDir(".", "")
	defer os.RemoveAll(tempDir)

	filePath := tempDir + "/" + "other.xml"
	rawOther := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ChangeSet>
			<RawSql>
				<Down>
					<Stmt>bar</Stmt>
				</Down>
			</RawSql>
		</ChangeSet>
	`
	if err := ioutil.WriteFile(filePath, []byte(rawOther), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	raw := `
		<ChangeLog tableName="table_name">
			<Variables>
				<Variable name="whoami">dbschema</Variable>
				<Variable name="whoareyou">idontknow</Variable>
			</Variables>

			<ChangeSet>
				<RawSql>
					<Up>
						<Stmt>foo</Stmt>
					</Up>
				</RawSql>
			</ChangeSet>

			<Import path="` + filePath + `" />

			<Variables>
				<Variable name="whoareyou">stranger</Variable>
			</Variables>
		</ChangeLog>
	`

	dec := xml.NewDecoder(strings.NewReader(raw))

	cl := &ChangeLog{}
	if err := dec.Decode(cl); err != nil {
		t.Fatal(err)
	}

	if *cl.TableName != "table_name" {
		t.Fatal()
	}
	if cl.LockTableName != nil {
		t.Fatal()
	}
	if len(cl.Variables.Values) != 3 {
		t.Fatal()
	}
	if len(cl.ChangeSets) != 2 {
		t.Fatal()
	}
}
