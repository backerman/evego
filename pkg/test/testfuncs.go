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

// Package test provides useful extensions to goconvey.
package test

import (
	"fmt"
	"strings"

	"github.com/backerman/evego/pkg/types"
)

// Component is an output of reprocessing.
// When reprocessed, an item returns one or more MaterialOutputs.
// We could use InventoryLine here, but the tests are much easier
// to write this way.
type Component struct {
	Name     string
	Quantity int
}

func (c Component) String() string {
	return fmt.Sprintf("[%vx %v]", c.Quantity, c.Name)
}

// ShouldHaveComposition compares a []component expected result
// against the actual []InventoryLine.
func ShouldHaveComposition(actual interface{}, expected ...interface{}) string {
	actualComps, ok := actual.(*[]types.InventoryLine)
	if !ok {
		return "Failed to cast actual to inventory line array"
	}
	expectedComps, ok := expected[0].([]Component)
	if !ok {
		return "Failed to cast expected to component array"
	}
	var messages []string
	if len(*actualComps) != len(expectedComps) {
		return fmt.Sprintf("Wrong number of components returned. Expected: %v; actual: %v", expectedComps, *actualComps)
	}
	for _, comp := range *actualComps {
		// this is hacky; fix
		myname := comp.Item.Name
		found := false
		for _, exp := range expectedComps {
			if exp.Name == myname {
				found = true
				if exp.Quantity != comp.Quantity {
					messages = append(messages,
						fmt.Sprintf("Expected %d %s; actually got %d", exp.Quantity, myname, comp.Quantity))
				}
			}
		}
		if !found {
			messages = append(messages, fmt.Sprintf("Spurious output %v", myname))
			break
		}
	}

	if len(messages) > 0 {
		return strings.Join(messages, "; ")
	}
	return ""
}
