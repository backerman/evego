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

package eveapi_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/types"
	. "github.com/smartystreets/goconvey/convey"
)

var testDbPath = "../../testdb.sqlite"

var testOutpostsXML = "../../testdata/test-outposts.xml"

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
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		x := eveapi.XMLAPI(ts.URL, db)

		Convey("Given a valid outpost ID", func() {
			outpostID := 61000854

			Convey("Its information is returned.", func() {
				expected := &types.Station{
					Name:            "4-EP12 VIII - 4-EP12 Inches for Mittens",
					ID:              outpostID,
					SystemID:        30004553,
					ConstellationID: 20000665,
					RegionID:        10000058,
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
