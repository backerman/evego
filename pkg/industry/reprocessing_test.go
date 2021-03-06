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

package industry_test

import (
	"testing"

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/industry"

	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	// Register SQLite3 and PgSQL drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var testDbDriver, testDbPath string

func init() {
	viper.SetDefault("DBDriver", "sqlite3")
	viper.SetDefault("DBPath", "../../testdb.sqlite")
	viper.SetEnvPrefix("EVEGO_TEST")
	viper.AutomaticEnv()
	testDbDriver = viper.GetString("DBDriver")
	testDbPath = viper.GetString("DBPath")
}

func TestReprocessingModules(t *testing.T) {
	Convey("Set up mock database", t, func() {
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		defer db.Close()

		Convey("Given a module", func() {
			gun, err := db.ItemForName("150mm Prototype Gauss Gun")
			So(err, ShouldBeNil)
			Convey("With reprocessing rate of 50% in sovnull, no tax, and no scrapmetal skills", func() {
				reproRate := 0.50
				quantity := 1
				taxRate := 0.00
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 5000000,
					ScrapmetalProcessing:   0,
					OreProcessing:          map[string]int{},
				}

				Convey("It should return the correct minerals", func() {
					reprocessed, err := industry.ReprocessItem(db, gun, quantity, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
						{"Tritanium", 614},
						{"Pyerite", 33},
						{"Mexallon", 38},
					})
				})
			})

			Convey("With reprocessing rate of 37.5% and level V scrapmetal skills", func() {
				reproRate := 0.375
				quantity := 1
				taxRate := 0.00
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 424,
					ScrapmetalProcessing:   5,
					OreProcessing:          map[string]int{},
				}

				Convey("It should return the correct minerals", func() {
					reprocessed, err := industry.ReprocessItem(db, gun, quantity, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
						{"Tritanium", 506},
						{"Pyerite", 27},
						{"Mexallon", 31},
					})
				})
			})

		})

		Convey("Given another module", func() {
			plate, err := db.ItemForName("800mm Crystalline Carbonide Restrained Plates")
			So(err, ShouldBeNil)

			Convey("With reprocessing rate of 50% and no scrapmetal skills", func() {
				reproRate := 0.50
				quantity := 1
				taxRate := 0.00
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 5000000,
					ScrapmetalProcessing:   0,
					OreProcessing:          map[string]int{},
				}

				Convey("It should return the correct minerals", func() {
					reprocessed, err := industry.ReprocessItem(db, plate, quantity, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
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
					reprocessed, err := industry.ReprocessItem(db, plate, quantity, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
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
				taxRate := 0.00
				skills := industry.ReproSkills{
					Reprocessing:           5000,
					ReprocessingEfficiency: 5000000,
					ScrapmetalProcessing:   0,
					OreProcessing:          map[string]int{},
				}

				Convey("Rounding should occur after summing all input units", func() {
					reprocessed, err := industry.ReprocessItem(db, gun, quantity, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
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
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		defer db.Close()

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

			items := []evego.InventoryLine{
				{Item: cscordite, Quantity: cscorditeQty},
				{Item: kernite, Quantity: kerniteQty},
				{Item: scordite, Quantity: scorditeQty},
			}

			Convey("In an NPC station with no standings", func() {
				reproRate := 0.5
				taxRate := 0.05 // 5% tax rate
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
					reprocessed, err := industry.ReprocessItem(db, cscordite, cscorditeQty, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
						{"Tritanium", 295459},
						{"Pyerite", 147729},
						{"Condensed Scordite", 25},
					})
				})

				Convey("One ore alone #2", func() {
					reprocessed, err := industry.ReprocessItem(db, kernite, kerniteQty, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
						{"Tritanium", 18173},
						{"Mexallon", 36218},
						{"Isogen", 18173},
						{"Luminous Kernite", 83},
					})
				})

				Convey("One ore alone #3", func() {
					reprocessed, err := industry.ReprocessItem(db, scordite, scorditeQty, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
						{"Tritanium", 83951},
						{"Pyerite", 41976},
						{"Scordite", 41},
					})
				})

				Convey("Three different ores.", func() {
					reprocessed, err := industry.ReprocessItems(db, items, reproRate, taxRate, skills)
					So(err, ShouldBeNil)
					So(reprocessed, ShouldHaveComposition, []Component{
						{"Tritanium", 397583},
						{"Pyerite", 189705},
						{"Mexallon", 36218},
						{"Isogen", 18173},
						{"Condensed Scordite", 25},
						{"Luminous Kernite", 83},
						{"Scordite", 41},
					})

				})

			})

		})

	})

}
