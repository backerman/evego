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
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testAssetsXML = "../../testdata/assets.xml"
)

func TestAccountAssets(t *testing.T) {
	Convey("Set up API interface", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testAssetsXML)
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
			charID := 94319654

			Convey("The character's assets are returned.", func() {
				expected := []evego.InventoryItem{
					{
						ItemID:        1016277603864,
						LocationID:    60000175,
						TypeID:        16273,
						Quantity:      400,
						Flag:          evego.InventoryFlag(4),
						Unpackaged:    false,
						BlueprintType: evego.NotBlueprint,
						Contents:      []evego.InventoryItem{},
					},
					{
						ItemID:        1016278183813,
						LocationID:    60000175,
						TypeID:        16273,
						Quantity:      400,
						Flag:          evego.InventoryFlag(4),
						Unpackaged:    false,
						BlueprintType: evego.NotBlueprint,
						Contents:      []evego.InventoryItem{},
					},
					{
						ItemID:        1018023320731,
						LocationID:    60002650,
						TypeID:        588,
						Quantity:      1,
						Flag:          evego.InventoryFlag(4),
						Unpackaged:    true,
						BlueprintType: evego.NotBlueprint,
						Contents: []evego.InventoryItem{
							{
								ItemID:        1018023320732,
								TypeID:        3636,
								Quantity:      1,
								Flag:          evego.InventoryFlag(27),
								Unpackaged:    true,
								BlueprintType: evego.NotBlueprint,
								Contents:      []evego.InventoryItem{},
							},
							{
								ItemID:        1018023320733,
								TypeID:        3651,
								Quantity:      1,
								Flag:          evego.InventoryFlag(28),
								Unpackaged:    true,
								BlueprintType: evego.NotBlueprint,
								Contents:      []evego.InventoryItem{},
							},
							{
								ItemID:        1018023320735,
								TypeID:        34,
								Quantity:      1,
								Flag:          evego.InventoryFlag(5),
								Unpackaged:    false,
								BlueprintType: evego.NotBlueprint,
								Contents:      []evego.InventoryItem{},
							},
						},
					},
					{
						ItemID:        1014948732937,
						LocationID:    60009514,
						TypeID:        33102,
						Quantity:      1,
						Flag:          evego.InventoryFlag(4),
						Unpackaged:    true,
						BlueprintType: evego.BlueprintCopy,
						Contents:      []evego.InventoryItem{},
					},
					{
						ItemID:        1019265170333,
						LocationID:    61000829,
						TypeID:        17366,
						Quantity:      1,
						Flag:          evego.InventoryFlag(4),
						Unpackaged:    true,
						BlueprintType: evego.NotBlueprint,
						Contents: []evego.InventoryItem{
							{
								ItemID:        1016396811813,
								TypeID:        27309,
								Quantity:      1,
								Flag:          evego.InventoryFlag(64),
								Unpackaged:    true,
								BlueprintType: evego.BlueprintCopy,
								Contents:      []evego.InventoryItem{},
							},
						},
					},
				}
				actual, err := x.Assets(key, charID)
				So(err, ShouldBeNil)

				expectedURL := fmt.Sprintf(
					"/char/AssetList.xml.aspx?characterID=%d&keyID=%d&vcode=%s",
					charID, key.KeyID, key.VerificationCode)
				So(actualURL, ShouldEqual, expectedURL)
				// expiry time minus "current time" is 6h
				expiration := cacheData.PutExpires
				So(cacheData.GetKeys, ShouldContainKey, ts.URL+expectedURL)
				So(cacheData.PutKeys, ShouldContainKey, ts.URL+expectedURL)
				now := time.Now()
				So(expiration, ShouldHappenAfter, now)
				So(expiration, ShouldHappenWithin, 361*time.Minute, now)
				So(expiration, ShouldNotHappenWithin, 359*time.Minute, now)
				So(actual, ShouldResemble, expected)
			})
		})
	})
}
