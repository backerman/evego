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

		Convey("With an invalid system name", func() {
			systemName := "Oniboshi"

			Convey("An error is returned.", func() {
				_, err := db.SolarSystemForName(systemName)
				So(err, ShouldNotBeNil)
			})
		})

	})

}
