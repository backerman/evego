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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/types"
	. "github.com/smartystreets/goconvey/convey"

	// Register SQLite3 driver for static database export
	_ "github.com/mattn/go-sqlite3"
)

var testDbPath = "../../testdb.sqlite"

var testOutpostsXML = "../../testdata/test-outposts.xml"
var testCharSheetXML = "../../testdata/test-charsheet.xml"

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
		x := eveapi.EveXMLAPI(ts.URL, db)

		Convey("Given a valid outpost ID", func() {
			outpostID := 61000854

			Convey("Its information is returned.", func() {
				expected := &types.Station{
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
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		x := eveapi.EveXMLAPI(ts.URL, db)

		Convey("Given a valid outpost name pattern", func() {
			outpostName := "%CAT%station"

			Convey("Matching outposts are returned.", func() {
				expected := &[]types.Station{
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
	})
}

func TestCharacterSheet(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testCharSheetXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		x := eveapi.EveXMLAPI(ts.URL, db)

		Convey("Given a character's API key", func() {
			characterID := 94319654
			keyID := 12345
			verificationCode := "abcdef12345"

			Convey("Its information is returned.", func() {
				expected := &types.CharacterSheet{
					Name:          "Arjun Kansene",
					ID:            characterID,
					Corporation:   "Center for Advanced Studies",
					CorporationID: 1000169,
					Skills: []types.Skill{
						{Name: "Gunnery", TypeID: 3300, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Small Hybrid Turret", TypeID: 3301, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Spaceship Command", TypeID: 3327, NumSkillpoints: 45255, Level: 4, Published: true},
						{Name: "Gallente Frigate", TypeID: 3328, NumSkillpoints: 512000, Level: 5, Published: true},
						{Name: "Mining", TypeID: 3386, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Mechanics", TypeID: 3392, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Science", TypeID: 3402, NumSkillpoints: 45255, Level: 4, Published: true},
						{Name: "Astrometrics", TypeID: 3412, NumSkillpoints: 135765, Level: 4, Published: true},
						{Name: "Power Grid Management", TypeID: 3413, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Hacking", TypeID: 21718, NumSkillpoints: 135765, Level: 4, Published: true},
					},
				}
				actual, expiration, err := x.CharacterSheet(characterID, keyID, verificationCode)
				So(err, ShouldBeNil)

				expectedURL := fmt.Sprintf(
					"/char/CharacterSheet.xml.aspx?characterID=%d&keyID=%d&vcode=%s",
					characterID, keyID, verificationCode)
				So(actualURL, ShouldEqual, expectedURL)
				// expiry time minus "current time" is 57 minutes
				So(expiration, ShouldHappenWithin, 58*time.Minute, time.Now())
				So(expiration, ShouldNotHappenWithin, 56*time.Minute, time.Now())
				So(actual, ShouldResemble, expected)
			})
		})

	})
}
