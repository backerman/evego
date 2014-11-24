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

// Package parsing extracts useful data from external input (EFT, copy-and-paste
// from the game client, etc.)
package parsing

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/types"
)

// ParseInventory extracts a item inventory copied from the EVE client.
// This can be from:
// * contract
// * ship/station/container inventory
// * personal assets view
func ParseInventory(pasted string, database dbaccess.EveDatabase) *[]types.InventoryLine {
	results := []types.InventoryLine{}
	// First pass - we're going to assume item in field 1, quantity field 2

	// Break into individual lines.
	reader := csv.NewReader(strings.NewReader(pasted))
	reader.Comma = '\t'
	lines, err := reader.ReadAll()
	if err != nil {
		// Unable to parse.
		return &results
	}
	for _, line := range lines {
		// Line is a []string with the fields in the line.
		name := line[0]
		var quantity int
		numScanned, _ := fmt.Sscanf(line[1], "%d", &quantity)
		if numScanned != 1 {
			// Couldn't scan a valid quantity.
			continue
		}
		// FIXME: If available, use group to disambiguate item.
		// group := line[2]
		// See if we have a match.
		item, err := database.ItemForName(name)
		if err != nil {
			// Didn't find a matching item.
			continue
		}
		results = append(results, types.InventoryLine{Item: *item, Quantity: quantity})
	}
	return &results
}
