#!/bin/sh
#
# Copyright © 2014–5 Brad Ackerman.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Parse command-line arguments.
set -- `getopt hd:m:e: "$@"`
if [ $# -lt 1 ]
then
  # getopt failed
  exit 66
fi
while [ $# -gt 0 ]
do
  case "$1" in
    -e) sqlite_exec="$2"; shift
    ;;
    -d) dbfile="$2"; shift
    ;;
    --) shift; break;;
    -m) spatialite_module="$2"; shift;;
    -*) echo >&2 \
"usage: $0 -d dbfile -e sqlite_executable -m spatialite_module"
        exit 42;;
    *) break;;
  esac
  shift
done

if [ -z "${dbfile}" ]
  then
  echo "missing parameter: -d dbfile.sqlite"
  exit 2
fi

if [ -z "${sqlite_exec}" ]
  then
  echo "missing parameter: -e sqlite_executable"
  exit 2
fi

if [ -z "${spatialite_module}" ]
  then
  echo "missing parameter: -m spatialite_module"
  exit 2
fi

# Make the database spatial; generate the jumps data
${sqlite_exec} ${dbfile} <<EOF
SELECT load_extension('${spatialite_module}');

-- Initialize Spatialite triggers, views, and tables
SELECT InitSpatialMetaData();

-- Solar systems map
SELECT AddGeometryColumn('mapSolarSystems', 'the_geom', -1, 'POINT', 'XYZ');
UPDATE mapsolarsystems SET the_geom=MakePointZ(x,y,z, -1);

-- Jumps

CREATE TABLE jumps_data (
  fromSystem integer NOT NULL,
  toSystem integer NOT NULL,
  cost integer NOT NULL DEFAULT 1,
  PRIMARY KEY (fromSystem, toSystem)
);
-- Has to be 2D because spatialite_network doesn't like 3D.
-- Not that the geometry is actually being used here, mind.
SELECT AddGeometryColumn('jumps_data', 'the_geom', -1, 'LINESTRING', 'XY');
INSERT INTO jumps_data(fromSystem, toSystem, cost, the_geom)
  SELECT fromSolarSystemID, toSolarSystemID, 1,
    MakeLine(CastToXY(fss.the_geom), CastToXY(tss.the_geom))
  FROM mapSolarSystemJumps j
  JOIN mapSolarSystems fss ON j.fromSolarSystemID = fss.solarSystemID
  JOIN mapSolarSystems tss ON j.toSolarSystemID = tss.solarSystemID;
EOF

spatialite_network -d ${dbfile} -T jumps_data -f fromSystem -t toSystem \
  -c cost -g the_geom --unidirectional -o jumps_net

# Set up the routing output.
${sqlite_exec} ${dbfile} <<EOF
SELECT load_extension('${spatialite_module}');
CREATE VIRTUAL TABLE jump_route using virtualnetwork(jumps_net);
EOF
