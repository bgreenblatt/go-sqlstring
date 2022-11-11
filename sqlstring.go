package sqlstring

import (
	"fmt"
	"strings"
)

type SQLString struct {
	buildsql strings.Builder
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
	quotedS := fmt.Sprintf("'%s'", s)
	t.buildsql.WriteString(quotedS)
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
	qSep := "'" + sep + "'"
	growSize := len(s) * (len(s[0]) + len(qSep))
	t.buildsql.Grow(growSize)
	t.buildsql.WriteString("'")
	t.buildsql.WriteString(strings.Join(s, qSep))
	t.buildsql.WriteString("'")
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

func NewSQLString() *SQLString {
	return &SQLString{}
}
