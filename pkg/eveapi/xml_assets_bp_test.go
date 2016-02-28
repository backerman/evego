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
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testBlueprintsXML = "../../testdata/blueprints.xml"
)

func TestBlueprints(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestPath := r.URL.Path
				var responsePath string
				switch requestPath {
				case "/char/AssetList.xml.aspx":
					responsePath = testAssetsXML
				case "/char/Blueprints.xml.aspx":
					responsePath = testBlueprintsXML
					actualURL = r.URL.String()
				default:
					log.Fatalf("Attempted to retrieve an unknown endpoint.")
				}
				respFile, err := os.Open(responsePath)
				c.So(err, ShouldBeNil)
				defer respFile.Close()
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
			charID := 94319654

			Convey("The character's blueprints are returned.", func() {
				expected := []evego.BlueprintItem{
					{
						ItemID:             1008668596403,
						LocationID:         61000500,
						StationID:          61000500,
						TypeID:             2509,
						TypeName:           "Nova Torpedo Blueprint",
						Flag:               evego.InventoryFlag(4),
						Quantity:           1,
						TimeEfficiency:     16,
						MaterialEfficiency: 10,
						NumRuns:            -1,
						IsOriginal:         true,
					},
					{
						ItemID:             1014948732937,
						LocationID:         60009514,
						StationID:          60009514,
						TypeID:             33102,
						TypeName:           "Medium Ancillary Armor Repairer Blueprint",
						Flag:               evego.InventoryFlag(4),
						Quantity:           1,
						TimeEfficiency:     0,
						MaterialEfficiency: 0,
						NumRuns:            25,
						IsOriginal:         false,
					},
					{
						ItemID:             1014949259398,
						LocationID:         60009514,
						StationID:          60009514,
						TypeID:             33077,
						TypeName:           "Small Ancillary Armor Repairer Blueprint",
						Flag:               evego.InventoryFlag(4),
						Quantity:           1,
						TimeEfficiency:     0,
						MaterialEfficiency: 0,
						NumRuns:            3,
						IsOriginal:         false,
					},
					{
						ItemID:             1016396811813,
						LocationID:         1019265170333,
						StationID:          61000829,
						TypeID:             27309,
						TypeName:           "Station Warehouse Container Blueprint",
						Flag:               evego.InventoryFlag(64),
						Quantity:           1,
						TimeEfficiency:     0,
						MaterialEfficiency: 7,
						NumRuns:            35,
						IsOriginal:         false,
					},
				}
				assets, err := x.Assets(key, charID)
				So(err, ShouldBeNil)
				actual, err := x.Blueprints(key, charID, assets)
				So(err, ShouldBeNil)

				expectedURL := fmt.Sprintf(
					"/char/Blueprints.xml.aspx?characterID=%d&keyID=%d&vcode=%s",
					charID, key.KeyID, key.VerificationCode)
				So(actualURL, ShouldEqual, expectedURL)
				// expiry time minus "current time" is 11.5h
				expiration := cacheData.PutExpires
				So(cacheData.GetKeys, ShouldContainKey, ts.URL+expectedURL)
				So(cacheData.PutKeys, ShouldContainKey, ts.URL+expectedURL)
				now := time.Now()
				So(expiration, ShouldHappenAfter, now)
				So(expiration, ShouldHappenWithin, 691*time.Minute, now)
				So(expiration, ShouldNotHappenWithin, 689*time.Minute, now)
				So(actual, ShouldResemble, expected)
			})
		})
	})
}
