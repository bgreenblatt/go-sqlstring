// SqlString

// Copyright ©2022 Bruce Greenblatt

// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package sqlstring

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestRawUpdate(t *testing.T) {
	stmt := NewSQLString(true)

	stmt.AddString("UPDATE t1 SET name = ", false)
	stmt.AddStringWithQuotes("Bruce", false)
	stmt.AddString(" WHERE position = ", false)
	stmt.AddStringWithQuotes("engineer", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestUpdate(t *testing.T) {
	stmt := NewSQLStringUpdate(true)

	stmt.AddTable("t1", false)
	stmt.AddColumnValue("name", "Bruce", true)
	stmt.AddColumnValue("position", "Engineer", true)
	stmt.AddWhere("position == 'engineer'", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestRawSelect(t *testing.T) {
	var stmt SQLString

	stmt.AddString("SELECT c1 FROM t2 WHERE c2 = ", false)
	stmt.AddStringWithQuotes("ID2", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestRawSelect2(t *testing.T) {
	stmt := NewSQLString(false)

	stmt.AddString("SELECT c1 FROM t2 WHERE c2 = ", false)
	stmt.AddStringWithQuotes("ID2", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestSelect(t *testing.T) {
	var stmt SQLStringSelect

	stmt.AddColumn("c1", false)
	stmt.AddTable("t2", false)
	stmt.AddWhere("c1 == \"ID2\"", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestSelect2(t *testing.T) {
	var stmt SQLStringSelect

	stmt.AddColumn("c1", false)
	stmt.AddColumn("c2", false)
	stmt.AddTable("t2", false)
	stmt.AddWhere("c2 == 'ID2'", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestSelectGroupBy(t *testing.T) {
	var stmt SQLStringSelect

	stmt.AddColumn("c1", false)
	stmt.AddColumn("c2", false)
	stmt.AddTable("t2", false)
	stmt.AddWhere("c2 == 'ID2'", false)
	stmt.AddGroupBy("c2", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		// t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
		fmt.Printf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestSelectTableAlias(t *testing.T) {
	var stmt SQLStringSelect

	stmt.AddColumn("t.c1", false)
	stmt.AddColumn("t.c2", false)
	stmt.AddTable("t2 as t", false)
	stmt.AddWhere("c2 == 'ID2'", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		// t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
		fmt.Printf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestSelectDistinct(t *testing.T) {
	var stmt SQLStringSelect

	stmt.AddColumn("c1", false)
	stmt.AddColumn("c2", false)
	stmt.AddTable("t2", false)
	stmt.AddWhere("c2 == 'ID2'", false)
	stmt.AddAllDistinctOption(Distinct)

	err := checkSQL(t, stmt.String())
	if err != nil {
		// t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
		fmt.Printf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestInsert(t *testing.T) {
	stmt := NewSQLStringInsert(true)

	stmt.AddTable("t1", false)
	stmt.AddColumnValue("name", "Bruce", true)
	stmt.AddColumnValue("position", "Engineer", true)
	stmt.AddColumnValue("salary", "100000", false)

	fmt.Printf("SQL Insert statement is %s\n", stmt.String())
	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestInsertSelect(t *testing.T) {
	stmt := NewSQLStringInsert(false)
	selectStmt := NewSQLStringSelect(false)

	stmt.AddTable("t1", false)
	selectStmt.AddTable("t2", false)
	selectStmt.AddColumn("c1", false)
	selectStmt.AddColumn("c2", false)
	selectStmt.AddColumn("c3", false)
	selectStmt.AddWhere("c2 = 'ID2'", false)
	stmt.AddSelect(selectStmt)

	fmt.Printf("SQL Insert statement is %s\n", stmt.String())
	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Found error %v parsing: %s\n", err, stmt.String())
	}
}

func TestLaunchSQLite3(t *testing.T) {
	defer os.Remove("t2.db")
	err := createDB("t2.db")
	if err != nil {
		t.Errorf("Error creating DB: %v", err)
	}
	cmd := exec.Command("sqlite3", "t2.db", ".dump")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Error running command: %v", err)
	}
	fmt.Printf("%s\n", out)
}

func createDB(dbName string) error {
	schema1 := "CREATE TABLE t1 (name text primary key, position text, salary integer)"
	schema2 := "CREATE TABLE t2 (c1 text primary key, c2 text, c3 integer)"
	cmd := exec.Command("sqlite3", dbName, schema1)
	_, err := cmd.CombinedOutput()
	if err == nil {
		cmd := exec.Command("sqlite3", dbName, schema2)
		_, err = cmd.CombinedOutput()
	}
	return err
}

func checkSQL(t *testing.T, query string) error {
	err := createDB("t2.db")
	if err != nil {
		t.Errorf("Error creating DB: %v", err)
	}
	defer os.Remove("t2.db")

	cmd := exec.Command("sqlite3", "t2.db", query)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", out)
	cmd = exec.Command("sqlite3", "t2.db", ".dump")
	outTemp, _ := cmd.CombinedOutput()
	fmt.Printf("%s\n", outTemp)
	return err
}

func TestInsertExec(t *testing.T) {
	stmt := NewSQLStringInsert(true)

	stmt.AddTable("t1", false)
	stmt.AddColumnValue("name", "Bruce", true)
	stmt.AddColumnValue("position", "Engineer", true)
	stmt.AddColumnValue("salary", "100000", false)

	fmt.Printf("SQL Insert statement is %s\n", stmt.String())
	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Error running command: %v", err)
	}
}

func TestDelete(t *testing.T) {
	stmt := NewSQLStringDelete(true)

	stmt.AddTable("t2", false)
	stmt.AddWhere("c2 = 'ID2'", false)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Error running command: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	stmt := NewSQLStringCreateTable(true)

	dv1 := DefaultValue{
		Value:     "CURRENT_TIMESTAMP",
		UseQuotes: false,
	}
	dv2 := DefaultValue{
		Value:     "engineer",
		UseQuotes: true,
	}
	stmt.AddTable("t3", false)
	stmt.AddRow("c1", "TEXT", true, nil)
	stmt.AddRow("c2", "INTEGER", false, &dv1)
	stmt.AddRow("c3", "TEXT", false, nil)
	stmt.AddRow("c4", "TEXT", false, &dv2)

	err := checkSQL(t, stmt.String())
	if err != nil {
		t.Errorf("Error running command: %v", err)
	}
}
