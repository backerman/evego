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

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/backerman/evego"
	"github.com/jmoiron/sqlx"
)

type sqlDb struct {
	dbType                        databaseType
	db                            *sqlx.DB
	compStatement                 *sqlx.Stmt
	itemInfoStatement             *sqlx.Stmt
	itemIDInfoStatement           *sqlx.Stmt
	catTreeFromItemStatement      *sqlx.Stmt
	systemInfoStatement           *sqlx.Stmt
	systemIDInfoStatement         *sqlx.Stmt
	regionInfoStatement           *sqlx.Stmt
	stationIDInfoStatement        *sqlx.Stmt
	stationNameInfoStatement      *sqlx.Stmt
	blueprintProducesStmt         *sqlx.Stmt
	inputMaterialsToBlueprintStmt *sqlx.Stmt
	blueprintProducedByStmt       *sqlx.Stmt
	matsForBPProductionStmt       *sqlx.Stmt
	countJumpsStmt                *sqlx.Stmt
}

// DatabaseType is the SQL vendor that we're using.
type databaseType int

// The possible databse types
const (
	Unknown databaseType = iota
	SQLite
	PostgreSQL
)

// SQLDatabase returns an EveDatabase object that can be used to access an SQL backend.
func SQLDatabase(driver, dataSource string) evego.Database {
	evedb := new(sqlDb)
	var err error
	evedb.db, err = sqlx.Connect(driver, dataSource)
	db := evedb.db // shortcut
	if err != nil {
		log.Fatalf("Unable to open item database (driver: %s, datasource: %s): %v",
			driver, dataSource, err)
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
		{&evedb.stationNameInfoStatement, stationNameInfo},
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

	// Routing is not standardized, so we need to treat these statements
	// specially.
	if strings.Index(driver, "sqlite3") != -1 {
		// This is SQLite.
		evedb.dbType = SQLite
		evedb.countJumpsStmt, err = db.Preparex(countJumpsSQLite)
		// If the statement preparation returned an error and it's
		// "no such module: virtualnetwork", we're trying to do something
		// that requires Spatialite, but the Spatialite module has not been
		// loaded. Don't complain now; only complain if the caller attempts
		// to actually use spatial functionality.
		if err != nil && strings.Index(err.Error(), "virtualnetwork") < 0 {
			log.Fatalf("Unable to prepare statement: %v", err)
		}
	} else {
		// Unknown database. PostgreSQL will be supported, but it isn't
		// right now.
		log.Fatalf("Unknown database driver %v", driver)
	}
	return evedb
}

// ItemForName returns a populated Item object for a given item title.
func (db *sqlDb) ItemForName(itemName string) (*evego.Item, error) {
	var err error
	object := evego.Item{}
	row := db.itemInfoStatement.QueryRowx(itemName)
	err = row.StructScan(&object)
	if err == sql.ErrNoRows {
		return nil, err
	}
	object.Type = db.itemType(&object)

	return &object, err
}

func (db *sqlDb) ItemForID(itemID int) (*evego.Item, error) {
	var err error
	object := evego.Item{}
	row := db.itemIDInfoStatement.QueryRowx(itemID)
	err = row.StructScan(&object)
	if err == sql.ErrNoRows {
		return nil, err
	}

	object.Type = db.itemType(&object)
	return &object, err
}

// itemComposition returns the composition of a named Eve item.
func (db *sqlDb) ItemComposition(itemID int) ([]evego.InventoryLine, error) {
	rows, err := db.compStatement.Query(itemID)
	if err != nil {
		log.Fatalf("Unable to execute composition query for item %d: %v", itemID, err)
	}
	defer rows.Close()

	var results []evego.InventoryLine
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
		results = append(results, evego.InventoryLine{Quantity: quantity, Item: item})
	}
	return results, nil
}

// MarketGroupForItem returns the parent groups of the market item.
func (db *sqlDb) MarketGroupForItem(item *evego.Item) (*evego.MarketGroup, error) {
	rows, err := db.catTreeFromItemStatement.Query(item.ID)
	// The query doesn't return ErrNoRows, so we'll check for that case
	// below.
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer rows.Close()
	var itemGroup *evego.MarketGroup
	var curLevel *evego.MarketGroup
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
		nextLevel := &evego.MarketGroup{
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
	curLevel.Parent = &evego.MarketGroup{
		ID:          parentID,
		Name:        parentName,
		Parent:      nil,
		Description: parentDescription,
	}

	return itemGroup, nil
}

// itemType returns the type of this item, as required for reprocessing yield
// calculation. It's either ore, ice, or other.
func (db *sqlDb) itemType(item *evego.Item) evego.ItemType {
	catTree, err := db.MarketGroupForItem(item)
	if err != nil {
		return evego.UnknownItemType
	}
	for cur := catTree; cur != nil; cur = cur.Parent {
		if cur.Name == "Ore" {
			return evego.Ore
		}
		if cur.Name == "Ice Ore" {
			return evego.Ice
		}
	}
	// Ret
	return evego.Other
}

func (db *sqlDb) Close() error {
	return db.db.Close()
}

func (db *sqlDb) SolarSystemForName(systemName string) (*evego.SolarSystem, error) {
	row := db.systemInfoStatement.QueryRowx(systemName)
	system := &evego.SolarSystem{}
	err := row.StructScan(system)
	return system, err
}

func (db *sqlDb) SolarSystemsForPattern(systemName string) ([]evego.SolarSystem, error) {
	rows, err := db.systemInfoStatement.Queryx(systemName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var systems []evego.SolarSystem
	for rows.Next() {
		system := evego.SolarSystem{}
		rows.StructScan(&system)
		systems = append(systems, system)
	}
	if len(systems) == 0 {
		err = sql.ErrNoRows
	}
	return systems, err
}

func (db *sqlDb) SolarSystemForID(systemID int) (*evego.SolarSystem, error) {
	row := db.systemIDInfoStatement.QueryRowx(systemID)
	system := &evego.SolarSystem{}
	err := row.StructScan(system)
	return system, err
}

func (db *sqlDb) RegionForName(regionName string) (*evego.Region, error) {
	row := db.regionInfoStatement.QueryRowx(regionName)
	region := &evego.Region{}
	err := row.StructScan(region)
	return region, err
}

func (db *sqlDb) StationForID(stationID int) (*evego.Station, error) {
	row := db.stationIDInfoStatement.QueryRowx(stationID)
	station := &evego.Station{}
	err := row.StructScan(station)
	return station, err
}

func (db *sqlDb) StationsForName(stationName string) ([]evego.Station, error) {
	rows, err := db.stationNameInfoStatement.Queryx(stationName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stations []evego.Station
	for rows.Next() {
		station := evego.Station{}
		rows.StructScan(&station)
		stations = append(stations, station)
	}
	if len(stations) == 0 {
		err = sql.ErrNoRows
	}
	return stations, err
}

func activityToTypeCode(activityStr string) evego.ActivityType {
	switch activityStr {
	case "Manufacturing":
		return evego.Manufacturing
	case "Researching Technology":
		return evego.ResearchingTechnology
	case "Researching Time Efficiency":
		return evego.ResearchingTE
	case "Researching Material Efficiency":
		return evego.ResearchingME
	case "Copying":
		return evego.Copying
	case "Duplicating":
		return evego.Duplicating
	case "Reverse Engineering":
		return evego.ReverseEngineering
	case "Invention":
		return evego.Invention
	}
	// Unknown
	return evego.None
}

func (db *sqlDb) blueprintQuery(stmt *sqlx.Stmt, query string) ([]evego.IndustryActivity, error) {
	rows, err := stmt.Queryx(query)
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer rows.Close()
	var results []evego.IndustryActivity
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
		rowActivity := evego.IndustryActivity{
			InputItem:      input,
			OutputItem:     output,
			OutputQuantity: row.OutputProductQty,
		}
		rowActivity.ActivityType = activityToTypeCode(row.ActivityName)
		results = append(results, rowActivity)
	}

	return results, nil
}

func (db *sqlDb) BlueprintOutputs(typeName string) ([]evego.IndustryActivity, error) {
	return db.blueprintQuery(db.blueprintProducesStmt, typeName)
}

func (db *sqlDb) BlueprintForProduct(typeName string) ([]evego.IndustryActivity, error) {
	return db.blueprintQuery(db.blueprintProducedByStmt, typeName)
}

func (db *sqlDb) BlueprintsUsingMaterial(typeName string) ([]evego.IndustryActivity, error) {

	return db.blueprintQuery(db.inputMaterialsToBlueprintStmt, typeName)
}

func (db *sqlDb) BlueprintProductionInputs(
	typeName string, outputTypeName string) ([]evego.InventoryLine, error) {
	rows, err := db.matsForBPProductionStmt.Queryx(typeName, outputTypeName)
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer rows.Close()
	var results []evego.InventoryLine
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
		result := evego.InventoryLine{
			Quantity: row.InputMaterialQty,
			Item:     inputMat,
		}
		results = append(results, result)
	}

	return results, nil
}
