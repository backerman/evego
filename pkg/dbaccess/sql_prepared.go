/*
Copyright © 2014 Brad Ackerman.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package dbaccess

// Our hand-crafted SQL statements.
var (
	materialComposition = `
  SELECT mt.typeID AS materialID, m.quantity AS quantity
  FROM invTypes t, invTypes mt, invTypeMaterials m
  WHERE t.typeID = ?
  AND t.typeID = m.typeID
  AND mt.typeID = m.materialTypeID
  `
	itemInfo = `
  SELECT t.typeID, t.typeName, t.portionSize, g.groupName, c.categoryName
  FROM invTypes t, invCategories c, invGroups g
  WHERE t.typeName = ? AND t.groupID = g.groupID
  AND   g.categoryID = c.categoryID
  `
	itemIDInfo = `
  SELECT t.typeID, t.typeName, t.portionSize, g.groupName, c.categoryName
  FROM invTypes t, invCategories c, invGroups g
  WHERE t.typeID = ? AND t.groupID = g.groupID
  AND   g.categoryID = c.categoryID
  `
	catTree = `
  WITH RECURSIVE
  parents(marketGroupID, parentGroupID) AS
  (
    SELECT marketGroupID, parentGroupID FROM invMarketGroups
    WHERE marketGroupID = (
      SELECT marketGroupID
      FROM invTypes i
      JOIN invMarketGroups m USING(marketGroupID)
      WHERE i.typeID = ?
      )
      UNION ALL
      SELECT mg.marketGroupID, mg.parentGroupID
      FROM invMarketGroups mg
      INNER JOIN parents p ON mg.marketGroupID=p.parentGroupID
      )
      SELECT p.marketGroupID, m1.marketGroupName, m1.description, p.parentGroupID, m2.marketGroupName, m2.description
      FROM parents p
      JOIN invMarketGroups m1 ON p.marketGroupID = m1.marketGroupID
      JOIN invMarketGroups m2 ON p.parentGroupID = m2.marketGroupID
      `

	systemInfo = `
      SELECT s.solarSystemName, s.solarSystemID, s.security,
      c.constellationName, c.constellationID, r.regionName, r.regionID
      FROM   mapSolarSystems s
      JOIN   mapConstellations c USING(constellationID)
      JOIN   mapRegions r USING(regionID)
      WHERE  s.solarSystemName = ?
      `

	systemIDInfo = `
      SELECT s.solarSystemName, s.solarSystemID, s.security,
      c.constellationName, c.constellationID, r.regionName, r.regionID
      FROM   mapSolarSystems s
      JOIN   mapConstellations c USING(constellationID)
      JOIN   mapRegions r USING(regionID)
      WHERE  s.solarSystemID = ?
      `

	regionInfo = `
      SELECT regionid, regionname
      FROM   mapregions
      WHERE  regionName = ?
      `

	stationIDInfo = `
      SELECT stationName, stationID, solarSystemID, constellationID, regionID
      FROM   staStations
      WHERE  stationID = ?
      `
)
