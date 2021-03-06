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
#
# This script generates a test database from the full SDE dump
# (converted by Fuzzysteve) that is sufficiently small to be commited to
# version control.

DBLOC=${HOME}/Downloads/sqlite-latest.sqlite
TESTDB=./testdb.sqlite

# Dump schema for tables we use
sqlite3 $DBLOC > out-schema.sql <<EOF
.schema industryActivityMaterials
.schema industryActivityProducts
.schema invTypes
.schema invTypeMaterials
.schema invMarketGroups
.schema invCategories
.schema invGroups
.schema invNames
.schema mapSolarSystems
.schema mapSolarSystemJumps
.schema mapConstellations
.schema mapRegions
.schema ramActivities
.schema staStations
.quit
EOF

# Item names that we use in our tests.
ITEMS=$(cat<<EOF
(
"150mm Prototype Gauss Gun",
"800mm Crystalline Carbonide Restrained Plates",
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
"Beta Hull Mod Expanded Cargo",
"Tripped Power Circuit",
"Multifrequency S",
"Armor Plates",
"Small Supplemental Barrier Emitter I",
"Type-D Restrained Expanded Cargo",
"Limited Kinetic Plating I",
"Small I-ax Enduring Remote Armor Repairer",
"Shielded Radar Backup Cluster I",
"Medium Shield Extender II",
"EMP M",
"Vexor Blueprint",
"Vexor",
"Ishtar Blueprint",
"Datacore - Gallentean Starship Engineering",
"Datacore - Mechanical Engineering",
"Datacore - Amarrian Starship Engineering",
"Gunnery",
"Small Hybrid Turret",
"Spaceship Command",
"Gallente Frigate",
"Mining",
"Mechanics",
"Science",
"Astrometrics",
"Power Grid Management",
"Hacking",
"Liquid Ozone",
"Reaper",
"Civilian Gatling Autocannon",
"Civilian Miner",
"Tritanium",
"Medium Ancillary Armor Repairer Blueprint",
"Station Container",
"Station Warehouse Container Blueprint"
)
EOF
)

# Items that we use as input to blueprints.
ITEMS_BPINPUT=$(cat<<EOF
(
"Structure Laboratory"
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
"4-EP12",
"Polaris",
"Polfaly",
"Polstodur",
"8WA-Z6",
"Gisleres",
-- Routing....
"31-MLU",
"A9D-R0",
"BMNV-P",
"BY-S36",
"X-M2LR",
"FD-MLJ",
"PF-346",
"Orvolle",
-- Outposts
"C-OK0R",
"V-SEE6"
)
EOF
)

# Stations that we use in our tests.
STATIONS=$(cat<<EOF
(
"Alentene VII - Moon 5 - Astral Mining Inc. Refinery",
"Cistuvaert V - Moon 12 - Center for Advanced Studies School",
"Gisleres V - Moon 8 - Chemal Tech Factory",
"Gisleres IV - Moon 6 - Roden Shipyards Warehouse",
"Junsoraert XI - Moon 9 - Roden Shipyards Factory",
"Ouelletta V - Moon 5 - Federal Navy Academy",
"Sortet V - Moon 1 - Federation Navy Assembly Plant",
"Quier IV - Moon 27 - Sisters of EVE Treasury"
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

# Yeah, this is crazy complex.
BPMATWITH=$(cat<<EOF
WITH matTypes AS (
-- Stuff our test types are composed of
SELECT m.materialTypeID
FROM invTypes t, invTypeMaterials m
WHERE t.typeName IN ${ITEMS}
AND t.typeID = m.typeID
), bpInOut AS (
-- Stuff that uses or is produced using blueprint input
SELECT ti.typeName inputBP, tyo.typeName outputProduct
FROM   industryActivityProducts iap
JOIN   invTypes ti USING(typeID)
JOIN   invTypes tyo ON iap.productTypeID = tyo.typeID
JOIN   industryActivityMaterials iam
ON     iam.typeID = ti.typeID
JOIN   invTypes tm
ON     iam.materialTypeID = tm.typeID
WHERE  tm.typeName IN ${ITEMS_BPINPUT}
), bpTypes AS (
SELECT t.typeID
FROM invTypes t, bpInOut
-- test input materials
WHERE t.typeName IN ${ITEMS_BPINPUT}
-- what BPs require them
OR    t.typeName = bpInOut.inputBP
-- what they produce
OR    t.typeName = bpInOut.outputProduct
)
EOF
)

# Dump item data we need
sqlite3 $DBLOC <<EOF
.mode insert invTypes
.output out-data.sql
${BPMATWITH}
SELECT * FROM invTypes
WHERE typeName IN ${ITEMS}
OR    typeID IN matTypes
OR    typeID IN bpTypes;

.mode insert invTypeMaterials
WITH bpMats AS (
  SELECT typeID from invTypes
  WHERE  typeName IN ${ITEMS_BPINPUT}
),   items  AS (
  SELECT typeID from invTypes
  WHERE  typeName IN ${ITEMS}
)
SELECT * FROM invTypeMaterials
WHERE typeID IN bpMats
OR    typeID IN items
OR    materialTypeID IN bpMats;

.mode insert invGroups
${BPMATWITH}
SELECT * FROM invGroups
WHERE groupID IN (
  SELECT groupID from invTypes
  WHERE typeName IN ${ITEMS}
  OR    typeID IN matTypes
  OR    typeID IN bpTypes
);

.mode insert invCategories
${BPMATWITH}
SELECT * FROM invCategories
WHERE categoryID IN (
  SELECT categoryID from invGroups
  WHERE groupID IN (
    SELECT groupID from invTypes
    WHERE typeName IN ${ITEMS}
    OR    typeID IN matTypes
    OR    typeID IN bpTypes
    )
);

.mode insert invMarketGroups
${BPMATWITH}
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
      OR    i.typeID IN matTypes
      OR    i.typeID IN bpTypes
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

.mode insert mapSolarSystemJumps
WITH systemIDs AS (
  SELECT solarSystemID
  FROM   mapSolarSystems
  WHERE  solarSystemName IN $SYSTEMS
)
SELECT * FROM mapSolarSystemJumps
WHERE fromSolarSystemID IN systemIDs
OR    toSolarSystemID IN systemIDs;

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

.mode insert invNames
SELECT *
FROM   invNames
WHERE  itemID IN (
  SELECT DISTINCT corporationID
  FROM   staStations
  WHERE  stationName IN $STATIONS
);

-- Industry tables
.mode insert ramActivities
SELECT *
FROM   ramActivities;

.mode insert industryActivityMaterials
${BPMATWITH}
, types AS (
  SELECT typeID from invTypes
  WHERE typeName IN ${ITEMS}
)
SELECT *
FROM   industryActivityMaterials
WHERE  typeID IN types
OR     typeID IN bpTypes
OR     materialTypeID IN types;

.mode insert industryActivityProducts
${BPMATWITH}
, types AS (
SELECT typeID from invTypes
WHERE typeName IN ${ITEMS}
)
SELECT *
FROM   industryActivityProducts
WHERE  typeID IN types OR productTypeID IN types
OR     typeID IN bpTypes;
EOF

rm -f ${TESTDB}
sqlite3 ${TESTDB} <<EOF
BEGIN TRANSACTION;
.read out-schema.sql
.read out-data.sql
COMMIT;
EOF
