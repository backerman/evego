## Loading the SDE

The migrations in this directory assume that the database starts with a restore
of the [Fuzzwork SDE dump][1]. If you don't have such a database already,
download the `postgres-latest.dmp.bz2` file and run:
```
dropdb --if-exists evetool && createdb evetool
pg_restore -d evetool -O postgres-latest.dmp
```

To load the data into a schema other than public (here, we use `sde`):

```
dropdb --if-exists evetool && createdb evetool
bzcat postgres-latest.dmp.bz2| pg_restore -O | sed '/^SET search_path/ d' |
  (echo "create schema sde; set search_path to sde;" && cat ) |
  psql -1 evetool
```

Then apply the SQL files in this directory in numerical order:
```
psql evetool < *.sql

# or if you're doing the schema thing
# pgrouting v2.0.0 has a bug that prevents installing in a schema other than
# public
psql evetool <<EOF
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgrouting;
EOF
(echo "set search_path to sde,public;" && cat *.sql) | psql evetool
```

(evetool is the default database name, but you can choose a different one.)

This program requires PostgreSQL 9.3 or later; if for some reason you need to
use an earlier version, you will need to replace the materialized view in
`pkg/routing/pgsql_routing.sql` with a regular table.

## Upgrading the SDE version

To update the local SDE copy, use the `update_sde.py` script, e.g.:

```
update_sde.py /tmp/latest.dmp.bz2 myschema | psql -h somewhere mydatabase
```

[1]: https://www.fuzzwork.co.uk/dump/
