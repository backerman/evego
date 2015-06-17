The migrations in this directory assume that the database starts with a restore
of the [Fuzzwork SDE dump][1]. If you don't have such a database already,
download the `postgres-latest.dmp.bz2` file and run:
```
dropdb --if-exists evetool && createdb evetool
pg_restore -d evetool -O postgres-latest.dmp
```

Then apply them in numerical order:
```
psql evetool < *.sql
```

(evetool is the default database name, but you can choose a different one.)

This program requires PostgreSQL 9.3 or later; if for some reason you need to
use an earlier version, you will need to replace the materialized view in
`pkg/routing/pgsql_routing.sql` with a regular table.

[1]: https://www.fuzzwork.co.uk/dump/
