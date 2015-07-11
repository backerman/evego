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
	"strconv"
	"strings"
	"time"

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
	postgres: `
		SELECT COUNT(*)
		FROM eve_findRoute(?, ?)
		`,
}

type sqlRouter struct {
	db           *sqlx.DB
	numJumpsStmt *sqlx.Stmt
	dialect      dbType
	cache        evego.Cache
}

// SQLRouter returns a router that uses topological data stored in a SQL
// database. Currently, SQLite (with Spatialite) and PostgreSQL (with pgrouting
// and PostGIS) are supported.
func SQLRouter(driver, dataSource string, aCache evego.Cache) evego.Router {
	db, err := sqlx.Connect(driver, dataSource)
	if err != nil {
		log.Fatalf("Unable to open routing database (driver: %s, datasource: %s): %v",
			driver, dataSource, err)
	}
	var dialect dbType
	if strings.Index(driver, "sqlite3") != -1 {
		dialect = sqlite
	} else if driver == "postgres" {
		dialect = postgres
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
	return &sqlRouter{db: db, dialect: dialect, numJumpsStmt: numJumpsStmt, cache: aCache}
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

// putCache adds a routing result to the cache. The expiration time is
// arbitrarily set to one day.
func (r *sqlRouter) putCache(fromSystemID, toSystemID, numJumps int) error {
	key := "numjumps:" + strconv.Itoa(fromSystemID) +
		":" + strconv.Itoa(toSystemID)
	val := []byte(strconv.Itoa(numJumps))
	return r.cache.Put(key, val, time.Now().Add(24*time.Hour))
}

// getCache finds a routing result in the cache. It returns the number of jumps
// (undefined if not found) and whether the result was contained in the cache.
func (r *sqlRouter) getCache(fromSystemID, toSystemID int) (int, bool) {
	key := "numjumps:" + strconv.Itoa(fromSystemID) +
		":" + strconv.Itoa(toSystemID)
	val, found := r.cache.Get(key)
	if found {
		// Convert cached from []byte to integer and return it.
		cachedAsInt, err := strconv.Atoi(string(val))
		if err != nil {
			return 0, false
		}
		return cachedAsInt, true
	}
	return 0, false
}

func (r *sqlRouter) NumJumpsID(fromSystemID, toSystemID int) (int, error) {
	// This function will be implemented differently depending on the
	// backend database.
	if fromSystemID == toSystemID {
		// These are the same system.
		return 0, nil
	}

	// Check for this result in the cache; if it's already there, return it.
	cached, found := r.getCache(fromSystemID, toSystemID)
	if found {
		return cached, nil
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
			r.putCache(fromSystemID, toSystemID, -1)
			return -1, nil
		}
		r.putCache(fromSystemID, toSystemID, numRows-1)
		return numRows - 1, nil
	case postgres:
		var numRows int
		err := r.numJumpsStmt.Get(&numRows, fromSystemID, toSystemID)
		if err != nil {
			return 0, err
		}
		// numRows does not have a header. The row count n is zero for an impossible
		// route (e.g. anything to Polaris), one for origin and destination in the
		// same system, and k jumps where k=n+1 if n>=2. Since we've already checked
		// for the same-system case, that doesn't apply here.
		if numRows == 0 {
			r.putCache(fromSystemID, toSystemID, -1)
			return -1, nil
		}
		r.putCache(fromSystemID, toSystemID, numRows-1)
		return numRows - 1, nil
	default:
		return -1, fmt.Errorf("Routing is not supported for this database type.")
	}
}

func (r *sqlRouter) Close() error {
	return r.db.Close()
}
