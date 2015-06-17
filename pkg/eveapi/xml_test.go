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

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/cache"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	// Register SQLite3 and PgSQL drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testOutpostsXML      = "../../testdata/test-outposts.xml"
	testCharSheetXML     = "../../testdata/test-charsheet.xml"
	testAccountCharsXML  = "../../testdata/acct-characters.xml"
	testCharStandingsXML = "../../testdata/char-standings.xml"
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
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		cacheData := test.CacheData{}
		x := eveapi.XML(ts.URL, db, test.Cache(&cacheData))

		Convey("Given a character's API key", func() {
			characterID := 94319654
			key := &evego.XMLKey{
				KeyID:            12345,
				VerificationCode: "abcdef12345",
			}

			Convey("Its information is returned.", func() {
				expected := &evego.CharacterSheet{
					Character: evego.Character{
						Name:          "Arjun Kansene",
						ID:            characterID,
						Corporation:   "Center for Advanced Studies",
						CorporationID: 1000169,
						Alliance:      "C C P Alliance",
						AllianceID:    434243723,
					},
					Skills: []evego.Skill{
						{Name: "Gunnery", Group: "Gunnery", TypeID: 3300, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Small Hybrid Turret", Group: "Gunnery", TypeID: 3301, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Spaceship Command", Group: "Spaceship Command", TypeID: 3327, NumSkillpoints: 45255, Level: 4, Published: true},
						{Name: "Gallente Frigate", Group: "Spaceship Command", TypeID: 3328, NumSkillpoints: 512000, Level: 5, Published: true},
						{Name: "Mining", Group: "Resource Processing", TypeID: 3386, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Mechanics", Group: "Armor", TypeID: 3392, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Science", Group: "Science", TypeID: 3402, NumSkillpoints: 45255, Level: 4, Published: true},
						{Name: "Astrometrics", Group: "Scanning", TypeID: 3412, NumSkillpoints: 135765, Level: 4, Published: true},
						{Name: "Power Grid Management", Group: "Engineering", TypeID: 3413, NumSkillpoints: 256000, Level: 5, Published: true},
						{Name: "Hacking", Group: "Scanning", TypeID: 21718, NumSkillpoints: 135765, Level: 4, Published: true},
					},
				}
				actual, err := x.CharacterSheet(key, characterID)
				So(err, ShouldBeNil)

				expectedURL := fmt.Sprintf(
					"/char/CharacterSheet.xml.aspx?characterID=%d&keyID=%d&vcode=%s",
					characterID, key.KeyID, key.VerificationCode)
				So(actualURL, ShouldEqual, expectedURL)
				So(cacheData.GetKey, ShouldEqual, ts.URL+expectedURL)
				So(cacheData.PutKey, ShouldEqual, ts.URL+expectedURL)
				expiration := cacheData.PutExpires
				// expiry time minus "current time" is 57 minutes
				now := time.Now()
				So(expiration, ShouldHappenAfter, now)
				So(expiration, ShouldHappenWithin, 58*time.Minute, now)
				So(expiration, ShouldNotHappenWithin, 56*time.Minute, now)
				So(actual, ShouldResemble, expected)
			})
		})

	})
}

func TestAccountCharacters(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testAccountCharsXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		cacheData := test.CacheData{}
		x := eveapi.XML(ts.URL, db, test.Cache(&cacheData))

		Convey("Given an account's API key", func() {
			key := &evego.XMLKey{
				KeyID:            12345,
				VerificationCode: "abcdef12345",
			}

			Convey("The available characters on that account are returned.", func() {
				expected := []evego.Character{
					{
						Name:          "Arjun Kansene",
						ID:            94319654,
						Corporation:   "Center for Advanced Studies",
						CorporationID: 1000169,
					},
					{
						Name:          "All reps on Cain",
						ID:            123456,
						Corporation:   "Yes, this is test data",
						CorporationID: 78910,
						Alliance:      "Some Alliance",
						AllianceID:    494949,
					},
				}
				actual, err := x.AccountCharacters(key)
				So(err, ShouldBeNil)

				expectedURL := fmt.Sprintf(
					"/account/Characters.xml.aspx?keyID=%d&vcode=%s",
					key.KeyID, key.VerificationCode)
				So(actualURL, ShouldEqual, expectedURL)
				// expiry time minus "current time" is 38m16s
				expiration := cacheData.PutExpires
				So(cacheData.GetKey, ShouldEqual, ts.URL+expectedURL)
				So(cacheData.PutKey, ShouldEqual, ts.URL+expectedURL)
				now := time.Now()
				So(expiration, ShouldHappenAfter, now)
				So(expiration, ShouldHappenWithin, 39*time.Minute, now)
				So(expiration, ShouldNotHappenWithin, 38*time.Minute, now)
				So(actual, ShouldResemble, expected)
			})
		})
	})
}

func TestCharacterStandings(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testCharStandingsXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()
		db := dbaccess.SQLDatabase(testDbDriver, testDbPath)
		cacheData := test.CacheData{}
		x := eveapi.XML(ts.URL, db, test.Cache(&cacheData))

		Convey("Given an account's API key and a character ID", func() {
			key := &evego.XMLKey{
				KeyID:            12345,
				VerificationCode: "abcdef12345",
			}
			characterID := 94319654

			Convey("That character's standings are returned.", func() {
				expected := []evego.Standing{
					{EntityType: evego.NPCAgent, ID: 3009145,
						Name: "Ostes Zoenceliris", Standing: 1.06},
					{EntityType: evego.NPCAgent, ID: 3009372,
						Name: "Pauren Aubyrasse", Standing: 1.84},
					{EntityType: evego.NPCAgent, ID: 3009381,
						Name: "Arnerore Rylerave", Standing: 0.52},
					{EntityType: evego.NPCCorporation, ID: 1000005,
						Name: "Hyasyoda Corporation", Standing: 0.86},
					{EntityType: evego.NPCCorporation, ID: 1000010,
						Name: "Kaalakiota Corporation", Standing: 1.06},
					{EntityType: evego.NPCCorporation, ID: 1000017,
						Name: "Nugoeihuvi Corporation", Standing: 0.61},
					{EntityType: evego.NPCFaction, ID: 500001,
						Name: "Caldari State", Standing: -0.27},
					{EntityType: evego.NPCFaction, ID: 500002,
						Name: "Minmatar Republic", Standing: 0.95},
					{EntityType: evego.NPCFaction, ID: 500003,
						Name: "Amarr Empire", Standing: -2.41},
					{EntityType: evego.NPCFaction, ID: 500004,
						Name: "Gallente Federation", Standing: 0.77},
				}
				actual, err := x.CharacterStandings(key, characterID)
				So(err, ShouldBeNil)

				expectedURL := fmt.Sprintf(
					"/char/Standings.xml.aspx?characterID=%d&keyID=%d&vcode=%s",
					characterID, key.KeyID, key.VerificationCode)
				So(actualURL, ShouldEqual, expectedURL)
				expiration := cacheData.PutExpires
				So(cacheData.GetKey, ShouldEqual, ts.URL+expectedURL)
				So(cacheData.PutKey, ShouldEqual, ts.URL+expectedURL)
				// expiry time minus "current time" is 2h53m49s
				now := time.Now()
				So(expiration, ShouldHappenAfter, now)
				So(expiration, ShouldHappenWithin, 2*time.Hour+54*time.Minute, now)
				So(expiration, ShouldNotHappenWithin, 2*time.Hour+53*time.Minute, now)
				So(actual, ShouldResemble, expected)
			})
		})

	})
}
