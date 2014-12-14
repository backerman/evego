#!/bin/sh
#
# Copyright © 2014 Brad Ackerman.
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
#
# This script generates a test database from the full SDE dump
# (converted by Fuzzysteve) that is sufficiently small to be commited to
# version control.

DBLOC=${HOME}/Downloads/sqlite-latest.sqlite
TESTDB=./testdb.sqlite

# Dump schema for tables we use
sqlite3 $DBLOC > out-schema.sql <<EOF
.schema invTypes
.schema invTypeMaterials
.schema invMarketGroups
.schema invCategories
.schema invGroups
.schema mapSolarSystems
.schema mapConstellations
.schema mapRegions
.schema staStations
.quit
EOF

# Item names that we use in our tests.
ITEMS=$(cat<<EOF
(
"150mm Prototype Gauss Gun",
"800mm Reinforced Crystalline Carbonide Plates I",
"650mm Medium Carbine Howitzer I",
"Tritanium",
"Pyerite",
"Mexallon",
"Zydrine",
"Isogen",
"Nocxium",
"Megacyte",
"Condensed Scordite",
"Scordite",
"Luminous Kernite",
"Medium Automated Structural Restoration",
"Large Asymmetric Remote Capacitor Transmitter",
"Tripped Power Circuit",
"Multifrequency S",
"Armor Plates",
"Small Supplemental Barrier Emitter I",
"Beta Hull Mod Reinforced Bulkheads",
"Limited Kinetic Plating I",
"Small I-ax Remote Armor Repairer",
"Shielded Radar Backup Cluster I",
"Medium Shield Extender II",
"EMP M"
)
EOF
)

# System names that we use in our tests.
SYSTEMS=$(cat<<EOF
(
"Poitot",
"Dodixie",
"RF-GGF",
"J100015",
"4-EP12"
)
EOF
)

# Stations that we use in our tests.
STATIONS=$(cat<<EOF
(
"Alentene VII - Moon 5 - Astral Mining Inc. Refinery",
"Cistuvaert V - Moon 12 - Center for Advanced Studies School",
"Gisleres V - Moon 8 - Chemal Tech Factory",
"Junsoraert XI - Moon 9 - Roden Shipyards Factory",
"Ouelletta V - Moon 5 - Federal Navy Academy"
)
EOF
)

# Regions that we use in our tests.
REGIONS=$(cat<<EOF
(
"Outer Ring",
"Verge Vendor"
)
EOF
)


# Dump item data we need
sqlite3 $DBLOC <<EOF
.mode insert invTypes
.output out-data.sql
SELECT * FROM invTypes
WHERE typeName IN ${ITEMS}
OR    typeID IN (
  SELECT m.materialTypeID
  FROM invTypes t, invTypeMaterials m
  WHERE t.typeName IN ${ITEMS}
  AND t.typeID = m.typeID
);

.mode insert invTypeMaterials
SELECT * FROM invTypeMaterials
WHERE typeID IN (
  SELECT typeID from invTypes
  WHERE typeName IN ${ITEMS}
);

.mode insert invGroups
SELECT * FROM invGroups
WHERE groupID IN (
  SELECT groupID from invTypes
  WHERE typeName IN ${ITEMS}
  OR    typeID IN (
    SELECT m.materialTypeID
    FROM invTypes t, invTypeMaterials m
    WHERE t.typeName IN ${ITEMS}
    AND t.typeID = m.typeID
  )
);

.mode insert invCategories
SELECT * FROM invCategories
WHERE categoryID IN (
  SELECT categoryID from invGroups
  WHERE groupID IN (
    SELECT groupID from invTypes
    WHERE typeName IN ${ITEMS}
    OR    typeID IN (
      SELECT m.materialTypeID AS typeID
      FROM invTypes t, invTypeMaterials m
      WHERE t.typeName IN ${ITEMS}
      AND t.typeID = m.typeID
    )
  )
);

.mode insert invMarketGroups
SELECT * FROM invMarketGroups
WHERE marketGroupID IN (
WITH RECURSIVE
  parents(marketGroupID, parentGroupID) AS
    (
    SELECT marketGroupID, parentGroupID FROM invMarketGroups
    WHERE marketGroupID IN (
      SELECT marketGroupID
      FROM invTypes i
      JOIN invMarketGroups m USING(marketGroupID)
      WHERE i.typeName IN $ITEMS
      OR    i.typeID IN (
        SELECT m.materialTypeID AS typeID
        FROM invTypes t, invTypeMaterials m
        WHERE t.typeName IN ${ITEMS}
        AND t.typeID = m.typeID
      )
    )
    UNION ALL
    SELECT mg.marketGroupID, mg.parentGroupID
    FROM invMarketGroups mg
    INNER JOIN parents p ON mg.marketGroupID=p.parentGroupID
  )
SELECT p.marketGroupID AS groupId FROM parents p
UNION
SELECT p.parentGroupID AS groupId FROM parents p
);

.mode insert mapSolarSystems
SELECT * FROM mapSolarSystems
WHERE solarSystemName IN $SYSTEMS;

.mode insert mapConstellations
SELECT * FROM mapConstellations
WHERE constellationID IN (
  SELECT DISTINCT constellationID
  FROM mapSolarSystems
  WHERE solarSystemName IN $SYSTEMS
);

.mode insert mapRegions
SELECT * FROM mapRegions
WHERE regionID IN (
  SELECT DISTINCT regionID
  FROM mapSolarSystems
  WHERE solarSystemName IN $SYSTEMS
) OR regionName IN $REGIONS;

.mode insert staStations
SELECT *
FROM   staStations
WHERE  stationName IN $STATIONS;

EOF

rm -f ${TESTDB}
sqlite3 ${TESTDB} <<EOF
BEGIN TRANSACTION;
.read out-schema.sql
.read out-data.sql
COMMIT;
EOF
