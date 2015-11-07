/*
Copyright © 2014–5 Brad Ackerman.

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

import "strings"

// Our hand-crafted SQL statements.
var (
	materialComposition = `
  SELECT mt."typeID" AS "materialID", m."quantity" AS "quantity"
  FROM "invTypes" t, "invTypes" mt, "invTypeMaterials" m
  WHERE t."typeID" = ?
  AND t."typeID" = m."typeID"
  AND mt."typeID" = m."materialTypeID"
  `
	itemInfo = `
  SELECT t."typeID", t."typeName", t."portionSize", g."groupID", g."groupName", c."categoryName"
  FROM "invTypes" t, "invCategories" c, "invGroups" g
  WHERE t."typeName" = ? AND t."groupID" = g."groupID"
  AND   g."categoryID" = c."categoryID"
  `
	itemIDInfo = `
  SELECT t."typeID", t."typeName", t."portionSize", g."groupID", g."groupName", c."categoryName"
  FROM "invTypes" t, "invCategories" c, "invGroups" g
  WHERE t."typeID" = ? AND t."groupID" = g."groupID"
  AND   g."categoryID" = c."categoryID"
  `
	catTree = `
  WITH RECURSIVE
  parents("marketGroupID", "parentGroupID") AS
  (
    SELECT "marketGroupID", "parentGroupID" FROM "invMarketGroups"
    WHERE "marketGroupID" = (
      SELECT "marketGroupID"
      FROM "invTypes" i
      JOIN "invMarketGroups" m USING("marketGroupID")
      WHERE i."typeID" = ?
      )
      UNION ALL
      SELECT mg."marketGroupID", mg."parentGroupID"
      FROM "invMarketGroups" mg
      INNER JOIN parents p ON mg."marketGroupID"=p."parentGroupID"
      )
      SELECT p."marketGroupID", m1."marketGroupName", m1."description", p."parentGroupID", m2."marketGroupName", m2."description"
      FROM parents p
      JOIN "invMarketGroups" m1 ON p."marketGroupID" = m1."marketGroupID"
      JOIN "invMarketGroups" m2 ON p."parentGroupID" = m2."marketGroupID"
      `

	systemInfo = `
      SELECT   s."solarSystemName", s."solarSystemID", s."security",
               c."constellationName", c."constellationID", r."regionName", r."regionID"
      FROM     "mapSolarSystems" s
      JOIN     "mapConstellations" c USING("constellationID")
      JOIN     "mapRegions" r ON r."regionID" = c."regionID"
      WHERE    LOWER(s."solarSystemName") LIKE LOWER(?)
			ORDER BY s."solarSystemName"
      `

	systemIDInfo = `
      SELECT s."solarSystemName", s."solarSystemID", s."security",
             c."constellationName", c."constellationID", r."regionName", r."regionID"
      FROM   "mapSolarSystems" s
      JOIN   "mapConstellations" c USING("constellationID")
      JOIN   "mapRegions" r ON r."regionID" = c."regionID"
      WHERE  s."solarSystemID" = ?
      `

	regionInfo = `
      SELECT "regionID", "regionName"
      FROM   "mapRegions"
      WHERE  "regionName" = ?
      `

	stationIDInfo = `
      SELECT "stationName", "stationID", "solarSystemID", "constellationID", "regionID",
						 "corporationID", "itemName" "corporationName", "reprocessingEfficiency"
      FROM   "staStations" s
			JOIN   "invNames" n ON n."itemID" = s."corporationID"
      WHERE  "stationID" = ?
      `
	stationNameInfo = `
		SELECT "stationName", "stationID", "solarSystemID", "constellationID", "regionID",
					 "corporationID", "itemName" "corporationName", "reprocessingEfficiency"
		FROM   "staStations" s
		JOIN   "invNames" n ON n."itemID" = s."corporationID"
		WHERE  LOWER("stationName") LIKE LOWER(?)
		ORDER BY "stationName"
		`

	blueprintBase = `
		SELECT ti."typeName" "inputItem", ram."activityName", tyo."typeName" "outputProduct",
		       iap."quantity" "outputProductQty"
		FROM   "industryActivityProducts" iap
		JOIN   "invTypes" ti USING("typeID")
		JOIN   "ramActivities" ram USING("activityID")
		JOIN   "invTypes" tyo ON iap."productTypeID" = tyo."typeID"
		WHERE  QUERYCOLUMN LIKE ?
		ORDER BY "inputItem", "outputProduct"
		`

	// What items can I produce with a blueprint?
	blueprintProduces = strings.Replace(blueprintBase, "QUERYCOLUMN", "ti.\"typeName\"", 1)

	// How can I produce a blueprint?
	blueprintProducedBy = strings.Replace(blueprintBase, "QUERYCOLUMN", "tyo.\"typeName\"", 1)

	// Extra stanzas for WHERE when querying on input materials
	inputMatsWhere = `
		JOIN   "industryActivityMaterials" iam
		ON     iam."typeID" = ti."typeID"
		JOIN   "invTypes" tm
		ON     iam."materialTypeID" = tm."typeID"
	`
	inputMaterialsToBlueprint = strings.Replace(
		strings.Replace(blueprintBase, "WHERE", inputMatsWhere+" WHERE ", 1),
		"QUERYCOLUMN", "tm.\"typeName\"", 1)

	// Given a blueprint, what items do I need to manufacture/invent with it?
	materialsForBlueprintProduction = `
		SELECT ti."typeName" "inputItem", "activityName", tm."typeName" "inputMaterial",
					 iam."quantity" "inputMaterialQty", tyo."typeName" "outputProduct",
					 iap."quantity" "outputProductQty"
		FROM   "industryActivityMaterials" iam
		JOIN   "invTypes" ti USING("typeID")
		JOIN   "invTypes" tm
		ON     iam."materialTypeID" = tm."typeID"
		JOIN   "ramActivities" USING("activityID")
		JOIN   "industryActivityProducts" iap
		ON     iap."typeID" = ti."typeID" AND iap."activityID"=iam."activityID"
		JOIN   "invTypes" tyo ON iap."productTypeID" = tyo."typeID"
		WHERE  ti."typeName" = ? AND tyo."typeName" = ?
		ORDER BY "inputItem", "outputProduct", "inputMaterial"
		`

	// What are the possible outputs from reprocessing an item?
	reprocessOutputsStmt = `
		SELECT t_mat."typeID"
		FROM   "invTypes" t_mat
		JOIN   "invTypeMaterials" tm ON tm."materialTypeID" = t_mat."typeID"
		JOIN   "invTypes" t_prod ON tm."typeID" = t_prod."typeID"
		WHERE  t_prod."marketGroupID" IS NOT NULL
		GROUP BY t_mat."typeID"
		ORDER BY t_mat."typeName"
		`
)
