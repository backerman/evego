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

package parsing_test

import (
	"io/ioutil"
	"testing"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/parsing"
	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

var testDbPath = "../../testdb.sqlite"

func TestInventoryCopyPaste(t *testing.T) {

	Convey("Given a copy-and-pasted inventory", t, func() {
		inventory, err := ioutil.ReadFile("../../testdata/test-inventory.txt")
		inventoryStr := string(inventory)
		So(err, ShouldBeNil)

		db := dbaccess.SQLiteDatabase(testDbPath)
		Convey("It is correctly parsed.", func() {
			parsed := parsing.ParseInventory(inventoryStr, db)
			So(parsed, ShouldHaveComposition, []Component{
				{"Medium Automated Structural Restoration", 2},
				{"Large Asymmetric Remote Capacitor Transmitter", 1},
				{"Tripped Power Circuit", 42},
				{"Multifrequency S", 1},
				{"Armor Plates", 18},
				{"Small Supplemental Barrier Emitter I", 1},
				{"Beta Hull Mod Reinforced Bulkheads", 1},
				{"Limited Kinetic Plating I", 1},
				{"Small I-ax Remote Armor Repairer", 1},
				{"Shielded Radar Backup Cluster I", 1},
			})
		})

	})

}