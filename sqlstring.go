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
	"strings"
)

type SQLString struct {
	buildsql        strings.Builder
	useDoubleQuotes bool
}

type SelectAllUniqueOption int

const (
	None SelectAllUniqueOption = iota
	All
	Unique
)

func (s SelectAllUniqueOption) String() string {
	switch s {
	case None:
		return ""
	case All:
		return "ALL "
	case Unique:
		return "UNIQUE "
	default:
		return ""
	}
}

type SQLStringSelect struct {
	allOrUnique SelectAllUniqueOption
	sqlString   SQLString
	columns     []string
	tables      []string
	where       string
	groupby     []string
}

type SQLStringUpdate struct {
	sqlString SQLString
	columns   []string
	values    []string
	isquoted  []bool
	table     string
	where     string
}

func (t SQLString) String() string {
	return t.buildsql.String()
}

func (t *SQLString) AddString(s string, addComma bool) {
	t.buildsql.WriteString(s)
	if addComma {
		t.buildsql.WriteString(",")
	}
}

func (t *SQLString) AddStringWithQuotes(s string, addComma bool) {
	var quoteString string
	if t.useDoubleQuotes {
		quoteString = "\""
	} else {
		quoteString = "'"
	}
	t.buildsql.WriteString(quoteString)
	t.buildsql.WriteString(s)
	t.buildsql.WriteString(quoteString)
	if addComma {
		t.buildsql.WriteString(",")
	}
}

func (t *SQLString) AddStrings(s []string, sep string, addComma bool) {
	t.buildsql.WriteString(strings.Join(s, sep))
	if addComma {
		t.buildsql.WriteString(",")
	}
}

func (t *SQLString) AddStringsWithQuotes(s []string, sep string, addComma bool) {
	var quoteString string
	if t.useDoubleQuotes {
		quoteString = "\""
	} else {
		quoteString = "'"
	}
	qSep := quoteString + sep + quoteString
	growSize := len(s) * (len(s[0]) + len(qSep))
	t.buildsql.Grow(growSize)
	t.buildsql.WriteString(quoteString)
	t.buildsql.WriteString(strings.Join(s, qSep))
	t.buildsql.WriteString(quoteString)
	if addComma {
		t.buildsql.WriteString(",")
	}
}

func (t *SQLString) AddInt(i int, addComma bool) {
	s := fmt.Sprintf(" %d", i)
	t.buildsql.WriteString(s)
	if addComma {
		t.buildsql.WriteString(",")
	}
}

func (t *SQLString) Reset() {
	t.buildsql.Reset()
}

func NewSQLString(useDoubleQuotes bool) *SQLString {
	var stmt SQLString
	if useDoubleQuotes {
		stmt.useDoubleQuotes = true
	}
	return &stmt
}

func NewSQLStringSelect(useDoubleQuotes bool) *SQLStringSelect {
	var stmt SQLStringSelect
	if useDoubleQuotes {
		stmt.sqlString.useDoubleQuotes = true
	}
	return &stmt
}

func (t *SQLStringSelect) AddColumn(c string, addComma bool) {
	t.columns = append(t.columns, c)
}

func (t *SQLStringSelect) AddTable(tbl string, addComma bool) {
	t.tables = append(t.tables, tbl)
}

func (t *SQLStringSelect) AddWhere(w string, addComma bool) {
	t.where = w
}

func (t *SQLStringSelect) AddGroupBy(gb string, addComma bool) {
	t.groupby = append(t.groupby, gb)
}

func (t *SQLStringSelect) AddAllUniqueOption(s SelectAllUniqueOption) {
	t.allOrUnique = s
}

func (t *SQLStringSelect) String() string {
	t.sqlString.Reset()
	t.sqlString.AddString("SELECT ", false)
	t.sqlString.AddString(t.allOrUnique.String(), false)
	columns := strings.Join(t.columns, ", ")
	t.sqlString.AddString(columns, false)
	tables := strings.Join(t.tables, ", ")
	t.sqlString.AddString(" FROM ", false)
	t.sqlString.AddString(tables, false)
	t.sqlString.AddString(" WHERE ", false)
	t.sqlString.AddString(t.where, false)
	if len(t.groupby) > 0 {
		gbs := strings.Join(t.groupby, ", ")
		t.sqlString.AddString(" GROUP BY ", false)
		t.sqlString.AddString(gbs, false)
	}
	return t.sqlString.String()
}

func (t *SQLStringSelect) Reset() {
	t.sqlString.Reset()
}

func NewSQLStringUpdate(useDoubleQuotes bool) *SQLStringUpdate {
	var stmt SQLStringUpdate
	if useDoubleQuotes {
		stmt.sqlString.useDoubleQuotes = true
	}
	return &stmt
}

func (t *SQLStringUpdate) AddColumnValue(c string, v string, q bool) {
	t.columns = append(t.columns, c)
	t.values = append(t.values, v)
	t.isquoted = append(t.isquoted, q)
}

func (t *SQLStringUpdate) AddTable(tbl string, addComma bool) {
	t.table = tbl
}

func (t *SQLStringUpdate) AddWhere(w string, addComma bool) {
	t.where = w
}

func (t *SQLStringUpdate) String() string {
	t.sqlString.Reset()
	t.sqlString.AddString("UPDATE ", false)
	t.sqlString.AddString(t.table, false)
	t.sqlString.AddString(" SET ", false)
	for i := range t.columns {
		t.sqlString.AddString(t.columns[i], false)
		t.sqlString.AddString(" = ", false)
		var addComma bool
		if i+1 == len(t.columns) {
			addComma = false
		} else {
			addComma = true
		}
		if t.isquoted[i] {
			t.sqlString.AddStringWithQuotes(t.values[i], addComma)
		} else {
			t.sqlString.AddString(t.values[i], addComma)
		}
		if i+1 < len(t.columns) {
			t.sqlString.AddString(" ", false)
		}
	}
	t.sqlString.AddString(" WHERE ", false)
	t.sqlString.AddString(t.where, false)
	return t.sqlString.String()
}

func (t *SQLStringUpdate) Reset() {
	t.sqlString.Reset()
}
