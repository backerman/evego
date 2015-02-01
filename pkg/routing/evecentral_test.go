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
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/backerman/evego/pkg/cache"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/routing"

	. "github.com/smartystreets/goconvey/convey"

	// Register SQLite3 driver
)

var (
	actualURL   string
	fromToRegex *regexp.Regexp
)

func init() {
	fromToRegex = regexp.MustCompile("from/([^/]+)/to/([^?#]+)")
}

const (
	testRoutePrefix = "../../testdata/evecentral-route-"
	testRouteSuffix = ".json"

	// System IDs

	OrvolleID = "30003830"
	RfGgfID   = "30003333"
	BmnvPID   = "30003327"
	Xm2lrID   = "30003278"
	PolarisID = "30000380"
)

func getFromToSystem(urlString string) (string, string) {
	args := fromToRegex.FindStringSubmatch(actualURL)
	fromSystem, toSystem := args[1], args[2]
	return fromSystem, toSystem
}

func TestEVECentralRouting(t *testing.T) {
	Convey("Create the router struct.", t, func(c C) {
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				fromSys, toSys := getFromToSystem(actualURL)
				var whichResponse string
				switch fromSys {
				case "Orvolle", OrvolleID:
					switch toSys {
					case "Orvolle", OrvolleID:
						whichResponse = "samesystem"
					case "Polaris", PolarisID:
						whichResponse = "unreachable"
					case "RF-GGF", RfGgfID:
						whichResponse = "orvolle-rf"
					}
				case "BMNV-P", BmnvPID:
					if toSys == "X-M2LR" || toSys == Xm2lrID {
						whichResponse = "adjacent"
					}
				}
				if whichResponse == "" {
					log.Fatalf("No response available for route from %v to %v",
						fromSys, toSys)
				}
				respFile, err := os.Open(testRoutePrefix + whichResponse + testRouteSuffix)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()

		router := routing.EveCentralRouter(ts.URL, cache.NilCache())
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)

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
