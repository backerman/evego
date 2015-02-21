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

// Package parsing extracts useful data from external input (EFT, copy-and-paste
// from the game client, etc.)
package parsing

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/types"
)

var (
	// Matches a line such as "1,367,093 x Tritanium".
	industryLine = regexp.MustCompile(`^\s*([.\d,]+)\s*x?\s*(.*?)\s*$`)
)

// removeNonNumeric removes the separators (comma and/or full stop)
// from the string representation of an integer.
func removeNonNumeric(s string) string {
	newStrBuf := []rune{}
	for _, r := range s {
		if r >= '0' && r <= '9' {
			newStrBuf = append(newStrBuf, r)
		}
	}
	return string(newStrBuf)
}

// ParseInventory extracts a item inventory copied from the EVE client.
// This can be from:
// * contract
// * ship/station/container inventory
// * personal assets view
// * industry tab of item info
func ParseInventory(pasted string, database dbaccess.EveDatabase) []types.InventoryLine {
	results := []types.InventoryLine{}
	// Break into individual lines.
	reader := csv.NewReader(strings.NewReader(pasted))
	reader.Comma = '\t'
	lines, err := reader.ReadAll()
	if err != nil {
		// Unable to parse.
		return results
	}
	for _, line := range lines {
		var (
			quantity int
			itemName string
		)
		// Line is a []string with the fields in the line.
		if len(line) == 1 {
			// "123,456,789x Something"
			matches := industryLine.FindStringSubmatch(line[0])
			if matches == nil {
				continue
			}
			var err error
			_, err = fmt.Sscanf(removeNonNumeric(matches[1]), "%d", &quantity)
			if err != nil {
				continue
			}
			itemName = matches[2]
		} else {
			itemName = line[0]
			numScanned, _ := fmt.Sscanf(line[1], "%d", &quantity)
			if numScanned != 1 {
				// Couldn't scan a valid quantity.
				continue
			}
		}
		// FIXME: If available, use group to disambiguate item.
		// group := line[2]
		// See if we have a match.
		item, err := database.ItemForName(itemName)
		if err != nil {
			// Didn't find a matching item.
			continue
		}
		results = append(results, types.InventoryLine{Item: item, Quantity: quantity})

	}
	return results
}
