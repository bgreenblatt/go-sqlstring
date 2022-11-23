# sqlstring
Helper library for building SQL query strings

Notes:

1. This library uses strings.Builder since it is much faster than concatenating strings with "+".
1. Initial version doesn't really know anything about SQL just a convenience for building strings
1. Test code uses `github.com/krasun/gosqlparser` for SQL statement validation
