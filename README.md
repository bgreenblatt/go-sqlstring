# sqlstring
Helper library for building SQL query strings in Go

Notes:

1. This library uses strings.Builder since it is much faster than concatenating strings with "+".
1. Current version includes support for the following basic SQL queries:
    * SELECT
    * UPDATE
    * INSERT
1. There is also support for convenience functions for building strings so that arbitrary SQL queries can be created
1. Test code uses the sqlite3 command line tool for SQL statement validation
