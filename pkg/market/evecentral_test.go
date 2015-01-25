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

package market_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/market"
	"github.com/backerman/evego/pkg/types"
	. "github.com/smartystreets/goconvey/convey"

	// Register SQLite3 driver for static database export
	_ "github.com/mattn/go-sqlite3"
)

var testDbPath = "../../testdb.sqlite"

var testMarketOrdersXML = "../../testdata/test-marketorders.xml"
var testOutpostOrdersXML = "../../testdata/test-outpostorders.xml"
var testOutpostsXML = "../../testdata/test-outposts.xml"

type testElement struct {
	match string
	name  string
}

func shouldMatchOrders(actual interface{}, expected ...interface{}) string {
	actualOrders, ok := actual.(*[]types.Order)
	if !ok {
		return "Failed to cast actual to order array"
	}
	expectedOrders, ok := expected[0].(*[]types.Order)
	if !ok {
		return "Failed to cast actual to order array"
	}
	if len(*actualOrders) != len(*expectedOrders) {
		return fmt.Sprintf("Expected %d orders; received %d",
			len(*expectedOrders), len(*actualOrders))
	}
	var messages []string // errors found

	for i := range *expectedOrders {
		e, a := (*expectedOrders)[i], (*actualOrders)[i]
		tests := []testElement{
			{ShouldEqual(a.Type, e.Type), "Type"},
			{ShouldEqual(a.Quantity, e.Quantity), "Quantity"},
			{ShouldEqual(a.Item.ID, e.Item.ID), "Item"},
			{ShouldAlmostEqual(a.Price, e.Price), "Price"},
			{ShouldEqual(a.Station.ID, e.Station.ID), "Station"},
			{ShouldEqual(a.Station.Name, e.Station.Name), "Station name"},
			{ShouldBeTrue(a.Expiration.Equal(e.Expiration)), "Expiration date"},
		}
		if e.Type == types.Buy {
			tests = append(tests, []testElement{
				{ShouldEqual(a.MinQuantity, e.MinQuantity), "Minimum quantity"},
				{ShouldEqual(a.JumpRange, e.JumpRange), "Order range"},
			}...)
			if e.JumpRange == types.BuyNumberJumps {
				tests = append(tests, testElement{ShouldEqual(a.NumJumps, e.NumJumps), "Number of jumps"})
			}
		}
		for _, t := range tests {
			if t.match != "" {
				messages = append(messages, fmt.Sprintf("%s of order #%d doesn't match: %s",
					t.name, i, t.match))
			}
		}
	}

	if len(messages) > 0 {
		return strings.Join(messages, "; ")
	}
	return ""
}

func TestMarketOrders(t *testing.T) {

	// Convey relies on call-stack magic to pass its context struct;
	// this is broken by goroutines, so we need to explicitly get
	// the context (by accepting it in our test function) and call
	// the So function on it directly.
	Convey("Set up test data.", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testMarketOrdersXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))

		defer ts.Close()
		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		// We don't need outpost information here, so we don't pass in a reference
		// to the EVE XML API.
		ec := market.EveCentral(db, nil, ts.URL)

		Convey("Given a valid region and item", func() {
			regionName := "Verge Vendor"
			region, err := db.RegionForName(regionName)
			So(err, ShouldBeNil)
			orderType := types.AllOrders
			item, err := db.ItemForName("Medium Shield Extender II")
			So(err, ShouldBeNil)

			Convey("Results should be successfully processed.", func() {
				fna, err := db.StationForID(60014704)
				So(err, ShouldBeNil)
				amr, err := db.StationForID(60009556)
				So(err, ShouldBeNil)
				cas, err := db.StationForID(60014719)
				So(err, ShouldBeNil)
				ctf, err := db.StationForID(60010840)
				So(err, ShouldBeNil)
				expected := &[]types.Order{
					{Type: types.Sell, Item: item, Quantity: 20, Price: 999997.74, Station: fna,
						Expiration: time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC)},
					{Type: types.Sell, Item: item, Quantity: 4, Price: 1500000, Station: amr,
						Expiration: time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC)},
					{Type: types.Sell, Item: item, Quantity: 24, Price: 508989.90, Station: cas,
						Expiration: time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC)},
					{Type: types.Buy, Item: item, Quantity: 57, Price: 277000.00, Station: cas,
						Expiration:  time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC),
						MinQuantity: 1, JumpRange: types.BuyRegion},
					{Type: types.Buy, Item: item, Quantity: 64, Price: 0.01, Station: ctf,
						Expiration:  time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC),
						MinQuantity: 1, JumpRange: types.BuyNumberJumps, NumJumps: 10},
					{Type: types.Buy, Item: item, Quantity: 42, Price: 123.45, Station: ctf,
						Expiration:  time.Date(2015, time.January, 22, 0, 0, 0, 0, time.UTC),
						MinQuantity: 17, JumpRange: types.BuyStation},
					{Type: types.Buy, Item: item, Quantity: 1000, Price: 60000, Station: cas,
						Expiration:  time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC),
						MinQuantity: 1, JumpRange: types.BuySystem},
				}

				// Generate expected URL. url.Parse() sorts the value keys before
				// generating the query string, so the same parameters will always
				// generate the same result independent of the order in which they
				// are passed to Set.
				urlParms := url.Values{}
				urlParms.Set("typeid", fmt.Sprintf("%d", item.ID))
				urlParms.Set("regionlimit", fmt.Sprintf("%d", region.ID))
				expectedURL := "/?" + urlParms.Encode()

				actual, err := ec.OrdersForItem(item, regionName, orderType)
				So(err, ShouldBeNil)
				So(actualURL, ShouldEqual, expectedURL)
				So(actual, shouldMatchOrders, expected)
			})

		})

		Convey("Given an item and a station", func() {
			item, err := db.ItemForName("Medium Shield Extender II")
			So(err, ShouldBeNil)
			station, err := db.StationForID(60010336) // Gisleres IV-6 Roden Whse
			Convey("The buy orders valid at that station are identified.", func() {
				chemalTech, err := db.StationForID(60010840) // Gisleres V-8 Chemal Tech Factory
				So(err, ShouldBeNil)
				cas, err := db.StationForID(60014719) // Cistuvaert V-12 CAS
				So(err, ShouldBeNil)
				expected := &[]types.Order{
					{Type: types.Buy, Item: item, Quantity: 57, Price: 277000.00, Station: cas,
						Expiration:  time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC),
						MinQuantity: 1, JumpRange: types.BuyRegion},
					{Type: types.Buy, Item: item, Quantity: 64, Price: 0.01, Station: chemalTech,
						Expiration:  time.Date(2015, time.March, 2, 0, 0, 0, 0, time.UTC),
						MinQuantity: 1, JumpRange: types.BuyNumberJumps, NumJumps: 10},
				}
				actual, err := ec.BuyInStation(item, station)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})

	})

}

func TestOutpostOrders(t *testing.T) {
	Convey("Set up test data.", t, func(c C) {
		var actualURL string
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				respFile, err := os.Open(testOutpostOrdersXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))
		defer ts.Close()

		tsXMLAPI := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				respFile, err := os.Open(testOutpostsXML)
				c.So(err, ShouldBeNil)
				responseBytes, err := ioutil.ReadAll(respFile)
				c.So(err, ShouldBeNil)
				responseBuf := bytes.NewBuffer(responseBytes)
				responseBuf.WriteTo(w)
			}))
		defer tsXMLAPI.Close()

		db := dbaccess.SQLDatabase("sqlite3", testDbPath)
		xmlAPI := eveapi.XMLAPI(tsXMLAPI.URL, db)
		ec := market.EveCentral(db, xmlAPI, ts.URL)

		Convey("Given a valid location and item", func() {
			systemName := "4-EP12"
			system, err := db.SolarSystemForName(systemName)
			So(err, ShouldBeNil)
			orderType := types.AllOrders
			item, err := db.ItemForName("EMP M")
			So(err, ShouldBeNil)

			Convey("Results should be successfully processed.", func() {
				ifm, err := xmlAPI.OutpostForID(61000854) // 4-EP12 Inches for Mittens
				So(err, ShouldBeNil)
				expected := &[]types.Order{
					{Type: types.Sell, Item: item, Quantity: 35847, Price: 84.87, Station: ifm,
						Expiration: time.Date(2015, time.March, 14, 0, 0, 0, 0, time.UTC)},
					{Type: types.Sell, Item: item, Quantity: 16543, Price: 83.87, Station: ifm,
						Expiration: time.Date(2015, time.March, 14, 0, 0, 0, 0, time.UTC)},
				}

				urlParms := url.Values{}
				urlParms.Set("typeid", fmt.Sprintf("%d", item.ID))
				urlParms.Set("usesystem", fmt.Sprintf("%d", system.ID))
				expectedURL := "/?" + urlParms.Encode()

				actual, err := ec.OrdersForItem(item, systemName, orderType)
				So(err, ShouldBeNil)
				So(actualURL, ShouldEqual, expectedURL)
				So(actual, shouldMatchOrders, expected)
			})

		})

	})

}
