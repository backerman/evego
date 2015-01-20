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
	"fmt"

	"github.com/backerman/evego/pkg/types"
)

func (db *sqlDb) NumJumps(fromSystem, toSystem *types.SolarSystem) (int, error) {
	// This function will be implemented differently depending on the
	// backend database.
	if fromSystem.ID == toSystem.ID {
		// These are the same system.
		return 0, nil
	}
	switch db.dbType {
	case SQLite:
		var numRows int
		err := db.countJumpsStmt.Get(&numRows, fromSystem.ID, toSystem.ID)
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
