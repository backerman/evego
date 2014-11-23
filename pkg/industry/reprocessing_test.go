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

package industry_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/industry"
	"github.com/backerman/evego/pkg/types"

	. "github.com/smartystreets/goconvey/convey"
)

var testDbPath = "../../testdb.sqlite"

// component is an output of reprocessing.
// When reprocessed, an item returns one or more MaterialOutputs.
// We could use InventoryLine here, but the tests are much easier
// to write this way.
type component struct {
	Name     string
	Quantity int
}

func (c component) String() string {
	return fmt.Sprintf("[%vx %v]", c.Quantity, c.Name)
}

func shouldHaveComposition(actual interface{}, expected ...interface{}) string {
	actualComps, ok := actual.(*[]types.InventoryLine)
	if !ok {
		return "Failed to cast actual to inventory line array"
	}
	expectedComps, ok := expected[0].([]component)
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

func TestReprocessingModules(t *testing.T) {
	Convey("Set up mock database", t, func() {
		db := dbaccess.SQLiteDatabase(testDbPath)

		Convey("Given a module", func() {
			gun, err := db.ItemForName("150mm Prototype Gauss Gun")
			So(err, ShouldBeNil)
			Convey("With reprocessing rate of 50% in sovnull and no scrapmetal skills", func() {
				reproRate := 0.50
				quantity := 1
				standing := 10.0 // no tax
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 5000000,
					ScrapmetalProcessing:   0,
					OreProcessing:          map[string]int{},
				}

				Convey("It should return the correct minerals", func() {
					reprocessed := industry.ReprocessItem(gun, quantity, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 614},
						{"Pyerite", 33},
						{"Mexallon", 38},
					})
				})
			})

			Convey("With reprocessing rate of 37.5% and level V scrapmetal skills", func() {
				reproRate := 0.375
				quantity := 1
				standing := 10.0 // no tax
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 424,
					ScrapmetalProcessing:   5,
					OreProcessing:          map[string]int{},
				}

				Convey("It should return the correct minerals", func() {
					reprocessed := industry.ReprocessItem(gun, quantity, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 506},
						{"Pyerite", 27},
						{"Mexallon", 31},
					})
				})
			})

		})

		Convey("Given another module", func() {
			plate, err := db.ItemForName("800mm Reinforced Crystalline Carbonide Plates I")
			So(err, ShouldBeNil)

			Convey("With reprocessing rate of 50% and no scrapmetal skills", func() {
				reproRate := 0.50
				quantity := 1
				standing := 10.0 // no tax
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 5000000,
					ScrapmetalProcessing:   0,
					OreProcessing:          map[string]int{},
				}

				Convey("It should return the correct minerals", func() {
					reprocessed := industry.ReprocessItem(plate, quantity, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 5498},
						{"Pyerite", 5217},
						{"Mexallon", 3762},
						{"Isogen", 104},
						{"Nocxium", 2},
						{"Megacyte", 1},
					})
				})

				Convey("Rounding should occur after summing all input units", func() {
					quantity = 2
					reprocessed := industry.ReprocessItem(plate, quantity, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 10996},
						{"Pyerite", 10435},
						{"Mexallon", 7524},
						{"Isogen", 208},
						{"Nocxium", 5},
						{"Megacyte", 3},
					})
				})

			})
		})

		Convey("Given another gun", func() {
			gun, err := db.ItemForName("650mm Medium Carbine Howitzer I")
			So(err, ShouldBeNil)

			Convey("With reprocessing rate of 50% and no scrapmetal skills", func() {
				reproRate := 0.50
				quantity := 5
				standing := 10.0
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 5000000,
					ScrapmetalProcessing:   0,
					OreProcessing:          map[string]int{},
				}

				Convey("Rounding should occur after summing all input units", func() {
					reprocessed := industry.ReprocessItem(gun, quantity, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 10225},
						{"Pyerite", 3132},
						{"Mexallon", 3180},
						{"Isogen", 15},
					})
				})

			})
		})
	})

}

func TestReprocessingOre(t *testing.T) {

	Convey("Set up mock database", t, func() {
		db := dbaccess.SQLiteDatabase(testDbPath)

		Convey("Given some ore", func() {
			cscordite, err := db.ItemForName("Condensed Scordite")
			So(err, ShouldBeNil)
			cscorditeQty := 129825

			kernite, err := db.ItemForName("Luminous Kernite")
			So(err, ShouldBeNil)
			kerniteQty := 21083

			scordite, err := db.ItemForName("Scordite")
			So(err, ShouldBeNil)
			scorditeQty := 38841

			items := &[]types.InventoryLine{
				{Item: *cscordite, Quantity: cscorditeQty},
				{Item: *kernite, Quantity: kerniteQty},
				{Item: *scordite, Quantity: scorditeQty},
			}

			Convey("In an NPC station with no standings", func() {
				reproRate := 0.5
				standing := 0.0
				skills := industry.ReproSkills{
					Reprocessing:           5,
					ReprocessingEfficiency: 3,
					ScrapmetalProcessing:   0,
					OreProcessing: map[string]int{
						"Scordite": 4,
						"Kernite":  3,
					},
				}

				Convey("One ore alone #1", func() {
					reprocessed := industry.ReprocessItem(cscordite, cscorditeQty, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 294646},
						{"Pyerite", 147729},
						{"Condensed Scordite", 25},
					})
				})

				Convey("One ore alone #2", func() {
					reprocessed := industry.ReprocessItem(kernite, kerniteQty, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 18044},
						{"Mexallon", 36218},
						{"Isogen", 18044},
						{"Luminous Kernite", 83},
					})
				})

				Convey("One ore alone #3", func() {
					reprocessed := industry.ReprocessItem(scordite, scorditeQty, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 83951},
						{"Pyerite", 41976},
						{"Scordite", 41},
					})
				})

				Convey("Three different ores.", func() {
					reprocessed := industry.ReprocessItems(items, reproRate, standing, skills)
					So(reprocessed, shouldHaveComposition, []component{
						{"Tritanium", 396641},
						{"Pyerite", 189705},
						{"Mexallon", 36218},
						{"Isogen", 18044},
						{"Condensed Scordite", 25},
						{"Luminous Kernite", 83},
						{"Scordite", 41},
					})

				})

			})

		})

	})

}
