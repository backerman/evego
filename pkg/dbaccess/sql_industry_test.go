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

package dbaccess_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/dbaccess"

	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"

	// Register SQLite3 and PgSQL drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func shouldMatchActivities(actual interface{}, expected ...interface{}) string {
	actualIA, ok := actual.([]evego.IndustryActivity)
	if !ok {
		return "Failed to cast actual to []evego.IndustryActivity"
	}
	expectedIA, ok := expected[0].([]evego.IndustryActivity)
	if !ok {
		return "Failed to cast expected to []evego.IndustryActivity"
	}

	if len(actualIA) != len(expectedIA) {
		return fmt.Sprintf("Lenth mismatch: expected %d activities; received %d",
			len(expectedIA), len(actualIA))
	}

	// We know that actual and expected are the same length. Verify that
	// they have equivalent items (in any order).
	var sActual, sExpected, msgs []string
	for i := range actualIA {
		sActual = append(sActual, (actualIA)[i].String())
		sExpected = append(sExpected, (expectedIA)[i].String())
	}
	sort.Sort(sort.StringSlice(sActual))
	sort.Sort(sort.StringSlice(sExpected))
	for i := range sActual {
		if sActual[i] != sExpected[i] {
			msgs = append(msgs, fmt.Sprintf("Expected %v, got %v", sExpected[i], sActual[i]))
		}
	}

	return strings.Join(msgs, "; ")
}

func TestBlueprints(t *testing.T) {

	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)

		Convey("With a valid input blueprint", func() {
			typeName := "Vexor Blueprint"
			inputType, err := db.ItemForName(typeName)
			So(err, ShouldBeNil)

			Convey("We get correct products.", func() {
				vexor, err := db.ItemForName("Vexor")
				So(err, ShouldBeNil)
				ishtarBlueprint, err := db.ItemForName("Ishtar Blueprint")
				So(err, ShouldBeNil)
				expected := []evego.IndustryActivity{
					{InputItem: inputType,
						ActivityType:   evego.Manufacturing,
						OutputItem:     vexor,
						OutputQuantity: 1,
					},
					{InputItem: inputType,
						ActivityType:   evego.Invention,
						OutputItem:     ishtarBlueprint,
						OutputQuantity: 1,
					},
				}
				actual, err := db.BlueprintOutputs(typeName)
				So(err, ShouldBeNil)
				So(actual, shouldMatchActivities, expected)
			})
		})

		Convey("With a valid input material", func() {
			typeName := "Structure Laboratory"

			Convey("We get correct blueprints.", func() {
				expectedNames := []string{
					// Old-skool outposts
					"Amarr Factory Outpost Platform",
					"Caldari Research Outpost Platform",
					"Gallente Administrative Outpost Platform",
					// Citadels!
					"Fortizar",
					"Keepstar",
					// Citadel components?
					"Medium Laboratory",
					"Large Administration Hub",
					"Large Laboratory",
					"X-Large Administration Hub",
					"X-Large Laboratory",
				}
				var expected []evego.IndustryActivity
				for _, platform := range expectedNames {
					inBP, err := db.ItemForName(platform + " Blueprint")
					So(err, ShouldBeNil)
					outPlatform, err := db.ItemForName(platform)
					So(err, ShouldBeNil)
					expected = append(expected, evego.IndustryActivity{
						InputItem:      inBP,
						ActivityType:   evego.Manufacturing,
						OutputQuantity: 1,
						OutputItem:     outPlatform,
					})
				}
				actual, err := db.BlueprintsUsingMaterial(typeName)
				So(err, ShouldBeNil)
				So(actual, shouldMatchActivities, expected)
			})
		})

		Convey("With a valid blueprint and output", func() {
			inBP := "Vexor Blueprint"
			outBP := "Ishtar Blueprint"
			Convey("We get correct material requirements.", func() {
				expected := []Component{
					{Quantity: 8, Name: "Datacore - Gallentean Starship Engineering"},
					{Quantity: 8, Name: "Datacore - Mechanical Engineering"},
				}
				actual, err := db.BlueprintProductionInputs(inBP, outBP)
				So(err, ShouldBeNil)
				So(actual, ShouldHaveComposition, expected)
			})
		})

		Convey("With a valid output", func() {
			desiredOutput := "Ishtar Blueprint"
			Convey("We get the blueprint required to product it.", func() {
				inBP, err := db.ItemForName("Vexor Blueprint")
				So(err, ShouldBeNil)
				outBP, err := db.ItemForName(desiredOutput)
				So(err, ShouldBeNil)
				expected := []evego.IndustryActivity{
					{
						InputItem:      inBP,
						ActivityType:   evego.Invention,
						OutputItem:     outBP,
						OutputQuantity: 1,
					},
				}
				actual, err := db.BlueprintForProduct(desiredOutput)
				So(err, ShouldBeNil)
				So(actual, shouldMatchActivities, expected)
			})
		})

	})
}

// shouldContainItem takes a slice of Items and a type ID, and passes if some
// item in the slice has the input type ID.
func shouldContainItem(actual interface{}, expected ...interface{}) string {
	actualItems, ok := actual.([]evego.Item)
	if !ok {
		return "Failed to cast actual to []evego.Item"
	}
	expectedItem, ok := expected[0].(int)
	if !ok {
		return "Failed to cast expected to int"
	}
	for _, i := range actualItems {
		if i.ID == expectedItem {
			return ""
		}
	}
	return fmt.Sprintf("The item with type ID %v was not found in the actual items.", expectedItem)
}

// shouldNotContainItem matches the inverse of shouldContainItem.
func shouldNotContainItem(actual interface{}, expected ...interface{}) string {
	_, ok := actual.([]evego.Item)
	if !ok {
		return "Failed to cast actual to []evego.Item"
	}
	expectedItem, ok := expected[0].(int)
	if !ok {
		return "Failed to cast expected to int"
	}
	result := shouldContainItem(actual, expected...)
	if strings.Contains(result, "was not found") {
		// We're good.
		return ""
	}
	if result == "" {
		// Found the item.
		return fmt.Sprintf("The item with type ID %v was found in the actual items.", expectedItem)
	}
	// Not the blank string, but not the error we want to see.
	return result
}

func TestReprocessOutputs(t *testing.T) {
	Convey("Open a database connection.", t, func() {
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)

		Convey("Get the list of reprocessing outputs", func() {
			outputs, err := db.ReprocessOutputMaterials()
			So(err, ShouldBeNil)
			So(outputs, ShouldNotBeEmpty)
			Convey("Passes basic checks", func() {
				So(outputs, shouldContainItem, 11399)  // Morphite
				So(outputs, shouldContainItem, 11558)  // Sustained Shield Emitter
				So(outputs, shouldNotContainItem, 626) // Vexor
			})
		})

	})
}
