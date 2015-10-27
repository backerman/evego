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

package eveapi_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/cache"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	// Register SQLite3 and PgSQL drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testOutpostsXML = "../../testdata/test-outposts.xml"
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

func TestOutpostID(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testOutpostsXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		x := eveapi.XML(ts.URL, db, cache.NilCache())

		Convey("Given a valid outpost ID", func() {
			outpostID := 61000854

			Convey("Its information is returned.", func() {
				expected := &evego.Station{
					Name:            "4-EP12 VIII - 4-EP12 Inches for Mittens",
					ID:              outpostID,
					SystemID:        30004553,
					ConstellationID: 20000665,
					RegionID:        10000058,
					Corporation:     "GoonWaffe",
					CorporationID:   667531913,
				}
				actual, err := x.OutpostForID(outpostID)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})

		Convey("Given an invalid outpost ID", func() {
			outpostID := 321

			Convey("An error is returned.", func() {
				_, err := x.OutpostForID(outpostID)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestOutpostName(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testOutpostsXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		x := eveapi.XML(ts.URL, db, cache.NilCache())

		Convey("Given a valid outpost name pattern", func() {
			outpostName := "%CAT%station"

			Convey("Matching outposts are returned.", func() {
				expected := []evego.Station{
					{
						Name:            "8WA-Z6 VIII - CAT IN STATION",
						ID:              61000189,
						SystemID:        30004760,
						ConstellationID: 20000696,
						RegionID:        10000060,
						Corporation:     "Northern Associates Holdings",
						CorporationID:   98008728,
					},
				}
				actual, err := x.OutpostsForName(outpostName)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})

		Convey("Given an invalid outpost name", func() {
			outpostName := "Forty-two"

			Convey("An error is returned.", func() {
				_, err := x.OutpostsForName(outpostName)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Outposts can be dumped.", func() {
			outposts := x.DumpOutposts()
			So(len(outposts), ShouldEqual, 4)
		})
	})
}
