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

package routing_test

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/routing"

	. "github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	// Register SQLite3 and PgSQL drivers
	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
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

var (
	registerDriver sync.Once
)

func TestSQLRouting(t *testing.T) {
	Convey("Open a database connection.", t, func() {
		var router evego.Router
		var db evego.Database
		switch testDbDriver {
		case "sqlite3":
			// Register a custom SQLite3 driver with the Spatialite extension.
			// Has to be wrapped in a Once because this is executed multiple
			// times by GoConvey.
			registerDriver.Do(func() {
				sql.Register("sqlite3_spatialite",
					&sqlite3.SQLiteDriver{
						Extensions: []string{
							SpatialiteModulePath(),
						},
					})
			})
			router = routing.SQLRouter("sqlite3_spatialite", testDbPath)
			db = dbaccess.SQLDatabase("sqlite3_spatialite", testDbPath)
		case "postgres":
			router = routing.SQLRouter(testDbDriver, testDbPath)
			db = dbaccess.SQLDatabase(testDbDriver, testDbPath)
		default:
			Println("The database under test does not yet support routing; skipping.")
			return
		}

		defer db.Close()
		defer router.Close()

		Convey("Given a start and end system", func() {
			startSys, err := db.SolarSystemForName("Orvolle")
			So(err, ShouldBeNil)
			endSys, err := db.SolarSystemForName("RF-GGF")
			So(err, ShouldBeNil)

			Convey("The path is calculated correctly.", func() {
				numJumps, err := router.NumJumps(startSys, endSys)
				So(err, ShouldBeNil)
				So(numJumps, ShouldEqual, 5)
			})
		})

		Convey("Given an adjacent start and end system", func() {
			startSys, err := db.SolarSystemForName("BMNV-P")
			So(err, ShouldBeNil)
			endSys, err := db.SolarSystemForName("X-M2LR")
			So(err, ShouldBeNil)

			Convey("The path is calculated correctly.", func() {
				numJumps, err := router.NumJumps(startSys, endSys)
				So(err, ShouldBeNil)
				So(numJumps, ShouldEqual, 1)
			})
		})

		Convey("Given a start and end system that are the same", func() {
			startSys, err := db.SolarSystemForName("Orvolle")
			So(err, ShouldBeNil)
			endSys, err := db.SolarSystemForName("Orvolle")
			So(err, ShouldBeNil)

			Convey("The path is calculated correctly.", func() {
				numJumps, err := router.NumJumps(startSys, endSys)
				So(err, ShouldBeNil)
				So(numJumps, ShouldEqual, 0)
			})
		})

		Convey("Given an end system that cannot be reached from the start", func() {
			startSys, err := db.SolarSystemForName("Orvolle")
			So(err, ShouldBeNil)
			endSys, err := db.SolarSystemForName("Polaris")
			So(err, ShouldBeNil)

			Convey("Unreachability is correctly indicated.", func() {
				numJumps, err := router.NumJumps(startSys, endSys)
				So(err, ShouldBeNil)
				So(numJumps, ShouldEqual, -1)
			})
		})

	})
}
