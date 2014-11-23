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
	"log"

	"github.com/backerman/evego/pkg/types"
	// Register SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type sqliteDb struct {
	db                       *sql.DB
	compStatement            *sql.Stmt
	itemInfoStatement        *sql.Stmt
	itemIDInfoStatement      *sql.Stmt
	catTreeFromItemStatement *sql.Stmt
}

// SQLiteDatabase returns an EveDatabase object that can be used to access a SQLite backend.
func SQLiteDatabase(path string) EveDatabase {
	evedb := new(sqliteDb)
	var err error
	evedb.db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Unable to open item database: %v", err)
	}

	// Prepare statements
	evedb.compStatement, err = evedb.db.Prepare(`
	SELECT mt.typeID AS materialID, m.quantity AS quantity
	FROM invTypes t, invTypes mt, invTypeMaterials m
	WHERE t.typeID = ?
	AND t.typeID = m.typeID
	AND mt.typeID = m.materialTypeID
	`)
	if err != nil {
		log.Fatalf("Unable to prepare statement: %v", err)
	}
	evedb.itemInfoStatement, err = evedb.db.Prepare(`
	SELECT t.typeID, t.typeName, t.portionSize, g.groupName, c.categoryName
	FROM invTypes t, invCategories c, invGroups g
	WHERE t.typeName = ? AND t.groupID = g.groupID
	AND   g.categoryID = c.categoryID
	`)
	if err != nil {
		log.Fatalf("Unable to prepare statement: %v", err)
	}

	evedb.itemIDInfoStatement, err = evedb.db.Prepare(`
	SELECT t.typeID, t.typeName, t.portionSize, g.groupName, c.categoryName
	FROM invTypes t, invCategories c, invGroups g
	WHERE t.typeID = ? AND t.groupID = g.groupID
	AND   g.categoryID = c.categoryID
	`)
	if err != nil {
		log.Fatalf("Unable to prepare statement: %v", err)
	}

	evedb.catTreeFromItemStatement, err = evedb.db.Prepare(`
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
		`)

	return evedb
}

// ItemForName returns a populated Item object for a given item title.
func (db *sqliteDb) ItemForName(itemName string) (*types.Item, error) {
	var err error
	object := types.Item{}
	row := db.itemInfoStatement.QueryRow(itemName)
	err = row.Scan(&object.ID, &object.Name, &object.BatchSize, &object.Group, &object.Category)
	if err == sql.ErrNoRows {
		return nil, err
	}
	object.Type = db.itemType(&object)
	object.Materials, err = db.itemComposition(object.ID)
	return &object, err
}

func (db *sqliteDb) ItemForID(itemID int) (*types.Item, error) {
	var err error
	object := types.Item{}
	row := db.itemIDInfoStatement.QueryRow(itemID)
	err = row.Scan(&object.ID, &object.Name, &object.BatchSize, &object.Group, &object.Category)
	if err == sql.ErrNoRows {
		return nil, err
	}

	object.Materials, err = db.itemComposition(object.ID)
	object.Type = db.itemType(&object)
	return &object, err
}

// itemComposition returns the composition of a named Eve item.
func (db *sqliteDb) itemComposition(itemID int) ([]types.InventoryLine, error) {
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
		results = append(results, types.InventoryLine{Quantity: quantity, Item: *item})
	}
	return results, nil
}

// MarketGroupForItem returns the parent groups of the market item.
func (db *sqliteDb) MarketGroupForItem(item *types.Item) (*types.MarketGroup, error) {
	rows, err := db.catTreeFromItemStatement.Query(item.ID)
	if err == sql.ErrNoRows {
		return nil, err
	}
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
	for rows.Next() {
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
func (db *sqliteDb) itemType(item *types.Item) types.ItemType {
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

func (db *sqliteDb) Close() error {
	return db.db.Close()
}
