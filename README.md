# sqlstring
Helper library for building SQL query strings in Go

Notes:

1. This library uses strings.Builder since it is much faster than concatenating strings with "+".
1. Current version includes support for the following basic SQL queries:
    * SELECT
    * UPDATE
    * INSERT
1. There is also support for convenience functions for building strings so that arbitrary SQL queries can be created
1. Test code uses `github.com/krasun/gosqlparser` for SQL statement validation, but some test cases would fail since
   it doesn't support the full SQL Syntax.
