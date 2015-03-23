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

// Package routing is a thing now.
package routing

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/backerman/evego"
	"github.com/jmoiron/sqlx"
)

type dbType int

const (
	unknown dbType = iota
	sqlite
	postgres
)

var numJumpsSQL = map[dbType]string{
	sqlite: `
    SELECT COUNT(*)
    FROM   jump_route
    WHERE  NodeFrom = ? AND NodeTo = ?
    `,
}

type sqlRouter struct {
	db           *sqlx.DB
	numJumpsStmt *sqlx.Stmt
	dialect      dbType
}

// SQLRouter reutrns a thingy.
func SQLRouter(driver, dataSource string) evego.Router {
	db, err := sqlx.Connect(driver, dataSource)
	if err != nil {
		log.Fatalf("Unable to open routing database (driver: %s, datasource: %s): %v",
			driver, dataSource, err)
	}
	var dialect dbType
	if strings.Index(driver, "sqlite3") != -1 {
		dialect = sqlite
	} else {
		dialect = unknown
	}
	sqlStmt, ok := numJumpsSQL[dialect]
	if !ok {
		log.Fatalf("SQL driver %s is unsupported for routing", driver)
	}
	numJumpsStmt, err := db.Preparex(db.Rebind(sqlStmt))
	if err != nil {
		log.Fatalf("Unable to prepare jump plan statement: %v", err)
	}
	return &sqlRouter{db: db, dialect: dialect, numJumpsStmt: numJumpsStmt}
}

func (r *sqlRouter) NumJumps(fromSystem, toSystem *evego.SolarSystem) (int, error) {
	if fromSystem == nil {
		return 0, errors.New("Starting system must be non-nil")
	}
	if toSystem == nil {
		return 0, errors.New("Ending system must be non-nil")
	}
	return r.NumJumpsID(fromSystem.ID, toSystem.ID)
}

func (r *sqlRouter) NumJumpsID(fromSystemID, toSystemID int) (int, error) {
	// This function will be implemented differently depending on the
	// backend database.
	if fromSystemID == toSystemID {
		// These are the same system.
		return 0, nil
	}
	switch r.dialect {
	case sqlite:
		var numRows int
		err := r.numJumpsStmt.Get(&numRows, fromSystemID, toSystemID)
		if err != nil {
			return 0, err
		}
		// numRows has a header and then one row for each jump in the route.
		// If there is no route, we get a header and nothing else.
		//
		// Therefore, if numRows-1 is 0, there is no route; otherwise, the
		// route contains numRows-1 jumps.
		if numRows == 1 {
			return -1, nil
		}
		return numRows - 1, nil
	default:
		return -1, fmt.Errorf("Routing is not supported for this database type.")
	}
}

func (r *sqlRouter) Close() error {
	return nil
}
