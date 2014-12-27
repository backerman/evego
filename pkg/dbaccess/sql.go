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

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/backerman/evego/pkg/types"
	"github.com/jmoiron/sqlx"
)

type sqlDb struct {
	db                            *sqlx.DB
	compStatement                 *sqlx.Stmt
	itemInfoStatement             *sqlx.Stmt
	itemIDInfoStatement           *sqlx.Stmt
	catTreeFromItemStatement      *sqlx.Stmt
	systemInfoStatement           *sqlx.Stmt
	systemIDInfoStatement         *sqlx.Stmt
	regionInfoStatement           *sqlx.Stmt
	stationIDInfoStatement        *sqlx.Stmt
	blueprintProducesStmt         *sqlx.Stmt
	inputMaterialsToBlueprintStmt *sqlx.Stmt
	blueprintProducedByStmt       *sqlx.Stmt
	matsForBPProductionStmt       *sqlx.Stmt
}

// SQLDatabase returns an EveDatabase object that can be used to access an SQL backend.
func SQLDatabase(driver, dataSource string) EveDatabase {
	evedb := new(sqlDb)
	var err error
	evedb.db, err = sqlx.Connect(driver, dataSource)
	db := evedb.db // shortcut
	if err != nil {
		log.Fatalf("Unable to open item database: %v", err)
	}

	// Prepare statements
	stmts := []struct {
		preparedStatement **sqlx.Stmt
		statementText     string
	}{
		// Pointer magic, stage 1: Pass the address of the pointer.
		{&evedb.compStatement, materialComposition},
		{&evedb.itemInfoStatement, itemInfo},
		{&evedb.itemIDInfoStatement, itemIDInfo},
		{&evedb.catTreeFromItemStatement, catTree},
		{&evedb.systemInfoStatement, systemInfo},
		{&evedb.systemIDInfoStatement, systemIDInfo},
		{&evedb.regionInfoStatement, regionInfo},
		{&evedb.stationIDInfoStatement, stationIDInfo},
		{&evedb.blueprintProducesStmt, blueprintProduces},
		{&evedb.inputMaterialsToBlueprintStmt, inputMaterialsToBlueprint},
		{&evedb.blueprintProducedByStmt, blueprintProducedBy},
		{&evedb.matsForBPProductionStmt, materialsForBlueprintProduction},
	}

	for _, s := range stmts {
		prepared, err := db.Preparex(db.Rebind(s.statementText))
		if err != nil {
			log.Fatalf("Unable to prepare statement: %v", err)
		}
		// Pointer magic, stage 2: Dereference the pointer to the pointer
		// and set it to point to the statement we just prepared.
		*s.preparedStatement = prepared
	}

	return evedb
}

// ItemForName returns a populated Item object for a given item title.
func (db *sqlDb) ItemForName(itemName string) (*types.Item, error) {
	var err error
	object := types.Item{}
	row := db.itemInfoStatement.QueryRowx(itemName)
	err = row.StructScan(&object)
	if err == sql.ErrNoRows {
		return nil, err
	}
	object.Type = db.itemType(&object)
	object.Materials, err = db.itemComposition(object.ID)

	return &object, err
}

func (db *sqlDb) ItemForID(itemID int) (*types.Item, error) {
	var err error
	object := types.Item{}
	row := db.itemIDInfoStatement.QueryRowx(itemID)
	err = row.StructScan(&object)
	if err == sql.ErrNoRows {
		return nil, err
	}

	object.Materials, err = db.itemComposition(object.ID)
	object.Type = db.itemType(&object)
	return &object, err
}

// itemComposition returns the composition of a named Eve item.
func (db *sqlDb) itemComposition(itemID int) ([]types.InventoryLine, error) {
	rows, err := db.compStatement.Query(itemID)
	if err != nil {
		log.Fatalf("Unable to execute composition query for item %d: %v", itemID, err)
	}
	defer rows.Close()

	var results []types.InventoryLine
	for rows.Next() {
		var (
			id       int
			quantity int
		)
		err = rows.Scan(&id, &quantity)
		if err != nil {
			log.Fatalf("Unable to execute query for item %d: %v", itemID, err)
		}
		item, err := db.ItemForID(id)
		if err != nil {
			log.Fatalf("Unable to execute query for item %d component %d: %v", itemID, id, err)
		}
		results = append(results, types.InventoryLine{Quantity: quantity, Item: item})
	}
	return results, nil
}

// MarketGroupForItem returns the parent groups of the market item.
func (db *sqlDb) MarketGroupForItem(item *types.Item) (*types.MarketGroup, error) {
	rows, err := db.catTreeFromItemStatement.Query(item.ID)
	// The query doesn't return ErrNoRows, so we'll check for that case
	// below.
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer rows.Close()
	var itemGroup *types.MarketGroup
	var curLevel *types.MarketGroup
	// The SQL query returns the market group hierarchy for the queried item,
	// beginning with the item's group and walking the group hierarchy until
	// the most broad group is found.
	var (
		groupID           int
		groupName         string
		description       string
		parentID          int
		parentName        string
		parentDescription string
	)
	// Set hasRows to false here and true if we get some rows.
	hasRows := false
	for rows.Next() {
		hasRows = true
		rows.Scan(&groupID, &groupName, &description, &parentID, &parentName, &parentDescription)
		nextLevel := &types.MarketGroup{
			ID:          groupID,
			Name:        groupName,
			Parent:      nil,
			Description: description,
		}

		if itemGroup == nil {
			// itemGroup is our return value - the item's immediate parent group.
			itemGroup = nextLevel
		}
		if curLevel != nil {
			// This level is the parent of a previous level.
			curLevel.Parent = nextLevel
		}
		curLevel = nextLevel
	}
	if !hasRows {
		// No rows when expected, so return an error.
		return nil, sql.ErrNoRows
	}
	// The last row's parentID and parentName are the first-level market category;
	// add them as the final parent.
	curLevel.Parent = &types.MarketGroup{
		ID:          parentID,
		Name:        parentName,
		Parent:      nil,
		Description: parentDescription,
	}

	return itemGroup, nil
}

// itemType returns the type of this item, as required for reprocessing yield
// calculation. It's either ore, ice, or other.
func (db *sqlDb) itemType(item *types.Item) types.ItemType {
	catTree, err := db.MarketGroupForItem(item)
	if err != nil {
		log.Fatalf("Unable to get item type of item %v: %v", *item, err)
	}
	for cur := catTree; cur != nil; cur = cur.Parent {
		if cur.Name == "Ore" {
			return types.Ore
		}
		if cur.Name == "Ice Ore" {
			return types.Ice
		}
	}
	// Ret
	return types.Other
}

func (db *sqlDb) Close() error {
	return db.db.Close()
}

func (db *sqlDb) SolarSystemForName(systemName string) (*types.SolarSystem, error) {
	row := db.systemInfoStatement.QueryRowx(systemName)
	system := &types.SolarSystem{}
	err := row.StructScan(system)
	return system, err
}

func (db *sqlDb) SolarSystemsForPattern(systemName string) (*[]types.SolarSystem, error) {
	rows, err := db.systemInfoStatement.Queryx(systemName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var systems []types.SolarSystem
	for rows.Next() {
		system := types.SolarSystem{}
		rows.StructScan(&system)
		systems = append(systems, system)
	}
	if len(systems) == 0 {
		err = sql.ErrNoRows
	}
	return &systems, err
}

func (db *sqlDb) SolarSystemForID(systemID int) (*types.SolarSystem, error) {
	row := db.systemIDInfoStatement.QueryRowx(systemID)
	system := &types.SolarSystem{}
	err := row.StructScan(system)
	return system, err
}

func (db *sqlDb) RegionForName(regionName string) (*types.Region, error) {
	row := db.regionInfoStatement.QueryRowx(regionName)
	region := &types.Region{}
	err := row.StructScan(region)
	return region, err
}

func (db *sqlDb) StationForID(stationID int) (*types.Station, error) {
	row := db.stationIDInfoStatement.QueryRowx(stationID)
	station := &types.Station{}
	err := row.StructScan(station)
	return station, err
}

func activityToTypeCode(activityStr string) types.ActivityType {
	switch activityStr {
	case "Manufacturing":
		return types.Manufacturing
	case "Researching Technology":
		return types.ResearchingTechnology
	case "Researching Time Efficiency":
		return types.ResearchingTE
	case "Researching Material Efficiency":
		return types.ResearchingME
	case "Copying":
		return types.Copying
	case "Duplicating":
		return types.Duplicating
	case "Reverse Engineering":
		return types.ReverseEngineering
	case "Invention":
		return types.Invention
	}
	// Unknown
	return types.None
}

func (db *sqlDb) blueprintQuery(stmt *sqlx.Stmt, query string) (*[]types.IndustryActivity, error) {
	rows, err := stmt.Queryx(query)
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer rows.Close()
	var results []types.IndustryActivity
	for rows.Next() {
		row := struct {
			InputItem        string `db:"inputItem"`
			ActivityName     string `db:"activityName"`
			OutputProduct    string `db:"outputProduct"`
			OutputProductQty int    `db:"outputProductQty"`
		}{}
		err = rows.StructScan(&row)
		if err != nil {
			return nil, fmt.Errorf("Error parsing returned industry query: %v", err)
		}
		input, err := db.ItemForName(row.InputItem)
		if err != nil {
			// This would indicate some major error in the database, but we like
			// checking for such errors even if it puts a dent in our unit-test
			// coverage numbers.
			return nil, fmt.Errorf("Couldn't find item %v with row: %#v (%d results so far)",
				row.InputItem, row, len(results))
		}
		output, err := db.ItemForName(row.OutputProduct)
		if err != nil {
			return nil, err
		}
		rowActivity := types.IndustryActivity{
			InputItem:      input,
			OutputItem:     output,
			OutputQuantity: row.OutputProductQty,
		}
		rowActivity.ActivityType = activityToTypeCode(row.ActivityName)
		results = append(results, rowActivity)
	}

	return &results, nil
}

func (db *sqlDb) BlueprintOutputs(typeName string) (*[]types.IndustryActivity, error) {
	return db.blueprintQuery(db.blueprintProducesStmt, typeName)
}

func (db *sqlDb) BlueprintForProduct(typeName string) (*[]types.IndustryActivity, error) {
	return db.blueprintQuery(db.blueprintProducedByStmt, typeName)
}

func (db *sqlDb) BlueprintsUsingMaterial(typeName string) (*[]types.IndustryActivity, error) {

	return db.blueprintQuery(db.inputMaterialsToBlueprintStmt, typeName)
}

func (db *sqlDb) BlueprintProductionInputs(
	typeName string, outputTypeName string) (*[]types.InventoryLine, error) {
	rows, err := db.matsForBPProductionStmt.Queryx(typeName, outputTypeName)
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer rows.Close()
	var results []types.InventoryLine
	for rows.Next() {
		row := struct {
			InputItem        string `db:"inputItem"`
			ActivityName     string `db:"activityName"`
			InputMaterial    string `db:"inputMaterial"`
			OutputProduct    string `db:"outputProduct"`
			InputMaterialQty int    `db:"inputMaterialQty"`
			OutputProductQty int    `db:"outputProductQty"`
			Consume          bool   `db:"consume"`
		}{}
		rows.StructScan(&row)

		// Process into an InventoryLine.
		inputMat, err := db.ItemForName(row.InputMaterial)
		if err != nil {
			log.Fatalf("Database inconsistency error: item %#v not available; %v",
				row.InputItem, err)
		}
		result := types.InventoryLine{
			Quantity: row.InputMaterialQty,
			Item:     inputMat,
		}
		results = append(results, result)
	}

	return &results, nil
}
