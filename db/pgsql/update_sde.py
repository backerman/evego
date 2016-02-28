#!/usr/bin/python
# Update the SDE schema, disabling triggers.

from __future__ import print_function

import argparse
import os
import re
import subprocess
import sys


# Use both RawDescriptionHelpFormatter and ArgumentDefaultsHelpFormatter for
# the help message.
class CustomFormatter(argparse.ArgumentDefaultsHelpFormatter,
                      argparse.RawDescriptionHelpFormatter):
    pass

parser = argparse.ArgumentParser(
    description='Update the SDE schema with a new dump.',
    formatter_class=CustomFormatter,
    epilog='''
Pipe the output of this program to the psql command, specifying the
database to connect to and any other required options, e.g.:

    update_sde.py /tmp/latest.dmp.bz2 myschema | psql mydatabase
''')
parser.add_argument('dumpfile',
                    help='The bzip2\'ed, PostgreSQL-format SDE dump file.')
parser.add_argument('schema', default='sde',
                    help='The database schema that the SDE is loaded in.',
                    nargs='?')
args = parser.parse_args()

# Truncate tables before we populate them.
copytable = re.compile("^COPY\s+(\S+)")

# Suppress attempts to set the search path.
searchpath = re.compile("^SET search_path")

bunzip = subprocess.Popen(["bzcat", args.dumpfile], bufsize=-1,
                          stdout=subprocess.PIPE)

# Get a pg_restore job
pgrestore = subprocess.Popen(["pg_restore", "-O", "-a"],
                             bufsize=-1, stdin=bunzip.stdout,
                             stdout=subprocess.PIPE)

# Output prequel.
print("SET search_path TO {0}, public;".format(args.schema))
print("BEGIN;")
print("SET CONSTRAINTS ALL DEFERRED;")

found_search_path = False
for line in pgrestore.stdout:
    # Check for start of table data
    match = copytable.match(line)
    if match:
        # Truncate the table before inserting data.
        tableName = match.group(1)
        print("DELETE FROM {0} WHERE 1=1;".format(tableName))
    else:
        # Check for SET search_path and suppress it.
        if not found_search_path and searchpath.match(line):
            found_search_path = True
            continue
    print(line, end="")

print("COMMIT;")
print("VACUUM FULL;")

# Now regenerate geometry.
mydir = os.path.dirname(os.path.abspath(__file__))
update_sql = os.path.join(mydir, "002-afterupdate.sql")
with open(update_sql) as f:
    for line in f:
        print(line)
