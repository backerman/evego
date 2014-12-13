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
	"testing"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/types"

	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

var testDbPath = "../../testdb.sqlite"

func TestItems(t *testing.T) {
	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		defer db.Close()

		Convey("With a valid item name", func() {
			itemName := "Medium Shield Extender II"
			expected := types.Item{
				Name:      itemName,
				ID:        3831,
				Type:      types.Other,
				Category:  "Module",
				Group:     "Shield Extender",
				BatchSize: 1,
			}

			Convey("The correct information is returned.", func() {
				actual, err := db.ItemForName(itemName)
				So(err, ShouldBeNil)
				So(actual.Name, ShouldEqual, expected.Name)
				So(actual.ID, ShouldEqual, expected.ID)
				So(actual.Type, ShouldEqual, expected.Type)
				So(actual.Category, ShouldEqual, expected.Category)
				So(actual.Group, ShouldEqual, expected.Group)
				So(actual.BatchSize, ShouldEqual, expected.BatchSize)
				So(&actual.Materials, ShouldHaveComposition, []Component{
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

			Convey("We get correct information.", nil)
			expected := types.SolarSystem{
				Name:            "Poitot",
				ID:              30003271,
				Security:        -0.019552,
				Constellation:   "Z-6NQ6",
				ConstellationID: 20000478,
				Region:          "Syndicate",
				RegionID:        10000041,
			}
			actual, err := db.SolarSystemForName(systemName)
			// Can't use ShouldResemble because of the float.
			So(err, ShouldBeNil)
			So(actual.Name, ShouldEqual, expected.Name)
			So(actual.ID, ShouldEqual, expected.ID)
			So(actual.Security, ShouldAlmostEqual, expected.Security)
			So(actual.Constellation, ShouldEqual, expected.Constellation)
			So(actual.ConstellationID, ShouldEqual, expected.ConstellationID)
			So(actual.Region, ShouldEqual, expected.Region)
			So(actual.RegionID, ShouldEqual, expected.RegionID)
		})

		Convey("With a valid system ID", func() {
			systemID := 30003333

			Convey("We get correct information.", func() {
				expected := types.SolarSystem{
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
				So(actual.Name, ShouldEqual, expected.Name)
				So(actual.ID, ShouldEqual, expected.ID)
				So(actual.Security, ShouldAlmostEqual, expected.Security)
				So(actual.Constellation, ShouldEqual, expected.Constellation)
				So(actual.ConstellationID, ShouldEqual, expected.ConstellationID)
				So(actual.Region, ShouldEqual, expected.Region)
				So(actual.RegionID, ShouldEqual, expected.RegionID)
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

	})

}

func TestStations(t *testing.T) {

	Convey("Open a database connection", t, func() {
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)

		Convey("With a valid station ID", func() {
			stationID := 60010312

			Convey("We get correct information.", nil)
			expected := &types.Station{
				Name:            "Junsoraert XI - Moon 9 - Roden Shipyards Factory",
				ID:              60010312,
				SystemID:        30003016,
				ConstellationID: 20000441,
				RegionID:        10000037,
			}
			actual, err := db.StationForID(stationID)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("With an invalid station ID", func() {
			stationID := 42

			Convey("An error is returned.", func() {
				_, err := db.StationForID(stationID)
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

			Convey("We get correct information.", nil)
			expected := &types.Region{
				Name: "Outer Ring",
				ID:   10000057,
			}
			actual, err := db.RegionForName(regionName)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
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
