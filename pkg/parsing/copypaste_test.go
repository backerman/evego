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

package parsing_test

import (
	"io/ioutil"
	"testing"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/parsing"
	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"

	// Register SQLite3 driver for static database export
	_ "github.com/mattn/go-sqlite3"
)

var testDbPath = "../../testdb.sqlite"

func TestInventoryCopyPaste(t *testing.T) {

	Convey("Given a copy-and-pasted inventory", t, func() {
		inventory, err := ioutil.ReadFile("../../testdata/test-inventory.txt")
		inventoryStr := string(inventory)
		So(err, ShouldBeNil)

		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		defer db.Close()

		Convey("It is correctly parsed.", func() {
			parsed := parsing.ParseInventory(inventoryStr, db)
			So(parsed, ShouldHaveComposition, []Component{
				{"Medium Automated Structural Restoration", 2},
				{"Large Asymmetric Remote Capacitor Transmitter", 1},
				{"Tripped Power Circuit", 42},
				{"Multifrequency S", 1},
				{"Armor Plates", 18},
				{"Small Supplemental Barrier Emitter I", 1},
				{"Type-D Restrained Expanded Cargo", 1},
				{"Limited Kinetic Plating I", 1},
				{"Small I-ax Remote Armor Repairer", 1},
				{"Shielded Radar Backup Cluster I", 1},
			})
		})
	})
}

func TestIndustryCopyPaste(t *testing.T) {

	Convey("Given a copy-and-pasted industry material list", t, func() {
		inventory, err := ioutil.ReadFile("../../testdata/test-industrytab.txt")
		inventoryStr := string(inventory)
		So(err, ShouldBeNil)

		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		defer db.Close()

		Convey("It is correctly parsed.", func() {
			parsed := parsing.ParseInventory(inventoryStr, db)
			So(parsed, ShouldHaveComposition, []Component{
				{"Tritanium", 1367093},
				{"Pyerite", 630827},
				{"Mexallon", 60890},
				{"Isogen", 9475},
				{"Nocxium", 2387},
				{"Zydrine", 1426},
				{"Megacyte", 359},
			})
		})
	})
}

func TestBadCopyPaste(t *testing.T) {
	Convey("Given completely malformed input", t, func() {
		inventoryStr := "fred"
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		defer db.Close()

		Convey("It returns an empty result.", func() {
			parsed := parsing.ParseInventory(inventoryStr, db)
			So(parsed, ShouldBeEmpty)
		})
	})
}
