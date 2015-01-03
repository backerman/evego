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

package dbaccess_test

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/types"

	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"

	// Register SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var testDbPath = "../../testdb.sqlite"

type testElement struct {
	match string
	name  string
}

// shouldMatchSystems is a custom matcher for *[]types.Solarsystem;
// we can't use ShouldResemble because of the float in that struct.
func shouldMatchSystems(actual interface{}, expected ...interface{}) string {
	actualSystems, ok := actual.(*[]types.SolarSystem)
	if !ok {
		return "Failed to cast actual to *[]types.SolarSystem"
	}
	expectedSystems, ok := expected[0].(*[]types.SolarSystem)
	if !ok {
		return "Failed to cast expected to *[]types.SolarSystem"
	}
	if len(*actualSystems) != len(*expectedSystems) {
		return fmt.Sprintf("Expected %d systems; received %d",
			len(*expectedSystems), len(*actualSystems))
	}
	var messages []string // errors found

	for i := range *expectedSystems {
		e, a := (*expectedSystems)[i], (*actualSystems)[i]
		tests := []testElement{
			{ShouldEqual(a.Name, e.Name), "Name"},
			{ShouldEqual(a.ID, e.ID), "ID"},
			{ShouldEqual(a.Constellation, e.Constellation), "Constellation"},
			{ShouldEqual(a.ConstellationID, e.ConstellationID), "Constellation ID"},
			{ShouldEqual(a.Region, e.Region), "Region"},
			{ShouldEqual(a.RegionID, e.RegionID), "Region ID"},
			{ShouldAlmostEqual(a.Security, e.Security), "Security level"},
		}
		for _, t := range tests {
			if t.match != "" {
				messages = append(messages, fmt.Sprintf("%s of system #%d doesn't match: %s",
					t.name, i, t.match))
			}
		}
	}

	if len(messages) > 0 {
		return strings.Join(messages, "; ")
	}
	return ""
}

// shouldMatchSystem is a convenience method for shouldMatchSystems, and takes
// a *types.SolarSystem instead of a *[]types.SolarSystem.
func shouldMatchSystem(actual interface{}, expected ...interface{}) string {
	actualSystem, ok := actual.(*types.SolarSystem)
	if !ok {
		return "Failed to cast actual to *types.SolarSystem"
	}
	expectedSystem, ok := expected[0].(*types.SolarSystem)
	if !ok {
		return "Failed to cast expected to *types.SolarSystem"
	}
	return shouldMatchSystems(&[]types.SolarSystem{*actualSystem},
		&[]types.SolarSystem{*expectedSystem})
}

func TestItems(t *testing.T) {
	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		defer db.Close()

		Convey("With a valid item name", func() {
			itemName := "Medium Shield Extender II"
			itemID := 3831
			expected := types.Item{
				Name:      itemName,
				ID:        itemID,
				Type:      types.Other,
				Category:  "Module",
				Group:     "Shield Extender",
				BatchSize: 1,
			}

			Convey("The correct information is returned.", func() {
				actual, err := db.ItemForName(itemName)
				So(err, ShouldBeNil)
				mats, err := db.ItemComposition(itemID)
				So(err, ShouldBeNil)
				So(actual.Name, ShouldEqual, expected.Name)
				So(actual.ID, ShouldEqual, expected.ID)
				So(actual.Type, ShouldEqual, expected.Type)
				So(actual.Category, ShouldEqual, expected.Category)
				So(actual.Group, ShouldEqual, expected.Group)
				So(actual.BatchSize, ShouldEqual, expected.BatchSize)
				So(mats, ShouldHaveComposition, []Component{
					{"Tritanium", 1890},
					{"Pyerite", 456},
					{"Mexallon", 179},
					{"Isogen", 6},
					{"Hydrogen Batteries", 6},
					{"Morphite", 5},
					{"Sustained Shield Emitter", 6},
				})
			})
		})

		Convey("With an invalid item", func() {
			itemName := "W76 Thermonuclear Device"

			Convey("An appropriate error is returned.", func() {
				_, err := db.ItemForName(itemName)
				So(err, ShouldEqual, sql.ErrNoRows)
			})
		})
	})
}

func TestSolarSystems(t *testing.T) {

	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)

		Convey("With a valid system name", func() {
			systemName := "Poitot"

			Convey("We get correct information.", func() {
				expected := &types.SolarSystem{
					Name:            "Poitot",
					ID:              30003271,
					Security:        -0.019552,
					Constellation:   "Z-6NQ6",
					ConstellationID: 20000478,
					Region:          "Syndicate",
					RegionID:        10000041,
				}
				actual, err := db.SolarSystemForName(systemName)
				So(err, ShouldBeNil)
				So(actual, shouldMatchSystem, expected)
			})
		})

		Convey("With a valid system ID", func() {
			systemID := 30003333

			Convey("We get correct information.", func() {
				expected := &types.SolarSystem{
					Name:            "RF-GGF",
					ID:              30003333,
					Security:        -0.246618,
					Constellation:   "49A-BZ",
					ConstellationID: 20000485,
					Region:          "Syndicate",
					RegionID:        10000041,
				}
				actual, err := db.SolarSystemForID(systemID)
				So(err, ShouldBeNil)
				So(actual, shouldMatchSystem, expected)
			})
		})

		Convey("With a valid system name pattern", func() {
			systemName := "Pol%"

			Convey("We get correct information.", func() {
				expected := &[]types.SolarSystem{
					{
						Name:            "Polaris",
						ID:              30000380,
						Security:        -0.000633,
						Constellation:   "9RW5-Z",
						ConstellationID: 20000054,
						Region:          "UUA-F4",
						RegionID:        10000004,
					},
					{
						Name:            "Polfaly",
						ID:              30005048,
						Security:        0.830126,
						Constellation:   "Nimedaz",
						ConstellationID: 20000738,
						Region:          "Kor-Azor",
						RegionID:        10000065,
					},
					{
						Name:            "Polstodur",
						ID:              30003434,
						Security:        0.836097,
						Constellation:   "Frar",
						ConstellationID: 20000501,
						Region:          "Metropolis",
						RegionID:        10000042,
					},
				}
				actual, err := db.SolarSystemsForPattern(systemName)
				So(err, ShouldBeNil)
				So(actual, shouldMatchSystems, expected)
			})
		})

		Convey("With an invalid system name", func() {
			systemName := "Oniboshi"

			Convey("An error is returned.", func() {
				_, err := db.SolarSystemForName(systemName)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("With an invalid system ID", func() {
			systemID := 12345

			Convey("An error is returned.", func() {
				_, err := db.SolarSystemForID(systemID)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("With an invalid system pattern", func() {
			systemName := "Onibo%"

			Convey("An error is returned.", func() {
				_, err := db.SolarSystemsForPattern(systemName)
				So(err, ShouldNotBeNil)
			})
		})

	})

}

func TestStations(t *testing.T) {

	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)

		Convey("With a valid station ID", func() {
			stationID := 60010312

			Convey("We get correct information.", func() {
				expected := &types.Station{
					Name:            "Junsoraert XI - Moon 9 - Roden Shipyards Factory",
					ID:              60010312,
					SystemID:        30003016,
					ConstellationID: 20000441,
					RegionID:        10000037,
					Corporation:     "Roden Shipyards",
					CorporationID:   1000102,
				}
				actual, err := db.StationForID(stationID)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})

		Convey("With an invalid station ID", func() {
			stationID := 42

			Convey("An error is returned.", func() {
				_, err := db.StationForID(stationID)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("With a valid station name pattern", func() {
			stationName := "%sisters%treas%"

			Convey("We get correct information.", func() {
				expected := &[]types.Station{{
					Name:            "Quier IV - Moon 27 - Sisters of EVE Treasury",
					ID:              60012655,
					SystemID:        30003037,
					ConstellationID: 20000444,
					RegionID:        10000037,
					Corporation:     "Sisters of EVE",
					CorporationID:   1000130,
				}}
				actual, err := db.StationsForName(stationName)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})

		Convey("With an invalid station name", func() {
			stationName := "Space Station Three"

			Convey("An error is returned.", func() {
				_, err := db.StationsForName(stationName)
				So(err, ShouldNotBeNil)
			})
		})

	})
}

func TestRegions(t *testing.T) {

	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)

		Convey("With a valid region name", func() {
			regionName := "Outer Ring"

			Convey("We get correct information.", func() {
				expected := &types.Region{
					Name: "Outer Ring",
					ID:   10000057,
				}
				actual, err := db.RegionForName(regionName)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)

			})
		})

		Convey("With an invalid region name", func() {
			regionName := "sudo make me a sandwich"

			Convey("An error is returned.", func() {
				_, err := db.RegionForName(regionName)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func shouldMatchActivities(actual interface{}, expected ...interface{}) string {
	actualIA, ok := actual.(*[]types.IndustryActivity)
	if !ok {
		return "Failed to cast actual to *[]types.IndustryActivity"
	}
	expectedIA, ok := expected[0].(*[]types.IndustryActivity)
	if !ok {
		return "Failed to cast expected to *[]types.IndustryActivity"
	}

	if len(*actualIA) != len(*expectedIA) {
		return fmt.Sprintf("Lenth mismatch: expected %d activities; received %d",
			len(*expectedIA), len(*actualIA))
	}

	// We know that actual and expected are the same length. Verify that
	// they have equivalent items (in any order).
	var sActual, sExpected, msgs []string
	for i := range *actualIA {
		sActual = append(sActual, (*actualIA)[i].String())
		sExpected = append(sExpected, (*expectedIA)[i].String())
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
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)

		Convey("With a valid input blueprint", func() {
			typeName := "Vexor Blueprint"
			inputType, err := db.ItemForName(typeName)
			So(err, ShouldBeNil)

			Convey("We get correct products.", func() {
				vexor, err := db.ItemForName("Vexor")
				So(err, ShouldBeNil)
				ishtarBlueprint, err := db.ItemForName("Ishtar Blueprint")
				So(err, ShouldBeNil)
				expected := &[]types.IndustryActivity{
					{InputItem: inputType,
						ActivityType:   types.Manufacturing,
						OutputItem:     vexor,
						OutputQuantity: 1,
					},
					{InputItem: inputType,
						ActivityType:   types.Invention,
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
			typeName := "Station Construction Parts"

			Convey("We get correct blueprints.", func() {
				expectedNames := []string{
					"Amarr Factory Outpost Platform",
					"Caldari Research Outpost Platform",
					"Gallente Administrative Outpost Platform",
					"Minmatar Service Outpost Platform",
				}
				var expected []types.IndustryActivity
				for _, platform := range expectedNames {
					inBP, err := db.ItemForName(platform + " Blueprint")
					So(err, ShouldBeNil)
					outPlatform, err := db.ItemForName(platform)
					So(err, ShouldBeNil)
					expected = append(expected, types.IndustryActivity{
						InputItem:      inBP,
						ActivityType:   types.Manufacturing,
						OutputQuantity: 1,
						OutputItem:     outPlatform,
					})
				}
				actual, err := db.BlueprintsUsingMaterial(typeName)
				So(err, ShouldBeNil)
				So(actual, shouldMatchActivities, &expected)
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
				expected := &[]types.IndustryActivity{
					{
						InputItem:      inBP,
						ActivityType:   types.Invention,
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
