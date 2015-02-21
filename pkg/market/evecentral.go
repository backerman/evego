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

package market

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/backerman/evego/pkg/cache"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/backerman/evego/pkg/routing"
	"github.com/backerman/evego/pkg/types"
)

type eveCentral struct {
	db        dbaccess.EveDatabase
	router    routing.EveRouter
	xmlAPI    eveapi.XMLAPI
	endpoint  *url.URL
	http      http.Client
	respCache cache.Cache
}

// EveCentralCached returns an interface to the EVE-Central API.
// It takes as input an EveDatabase object and an HTTP endpoint;
// the latter should be http://api.eve-central.com/api/quicklook
// for the production EVE-Central instance.
func EveCentralCached(db dbaccess.EveDatabase, router routing.EveRouter, xmlAPI eveapi.XMLAPI, endpoint string,
	aCache cache.Cache) EveMarket {
	epURL, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("Invalid URL %v passed for Eve-Central endpoint: %v", endpoint, err)
	}
	ec := eveCentral{db: db, router: router, endpoint: epURL, xmlAPI: xmlAPI, respCache: aCache}
	return &ec
}

// EveCentral returns an uncached interface to the EVE-Central API.
// This should only be used if the caller will be handling caching.
func EveCentral(db dbaccess.EveDatabase, router routing.EveRouter, xmlAPI eveapi.XMLAPI, endpoint string) EveMarket {
	myCache := cache.NilCache()
	return EveCentralCached(db, router, xmlAPI, endpoint, myCache)
}

func (e *eveCentral) getURL(u string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "evego (https://github.com/backerman/evego)")
	resp, err := e.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

type order struct {
	RegionID          int     `xml:"region"`
	StationID         int     `xml:"station"`
	StationName       string  `xml:"station_name"`
	Security          float64 `xml:"security"`
	Range             int     `xml:"range"`
	Price             float64 `xml:"price"`
	QuantityAvailable int     `xml:"vol_remain"`
	MinimumVolume     int     `xml:"min_volume"`
	ExpirationDate    string  `xml:"expires"`
	ReportedTime      string  `xml:"reported_time"`
}

type quicklook struct {
	SellOrders []order `xml:"quicklook>sell_orders>order"`
	BuyOrders  []order `xml:"quicklook>buy_orders>order"`
}

func (e *eveCentral) processOrders(data *quicklook, item *types.Item, t types.OrderType) []types.Order {
	var toProcess *[]order
	stationCache := make(map[int]*types.Station)
	switch t {
	case types.Buy:
		toProcess = &data.BuyOrders
	case types.Sell:
		toProcess = &data.SellOrders
	}
	results := []types.Order{}
	for _, o := range *toProcess {
		if stationCache[o.StationID] == nil {
			sta, err := e.db.StationForID(o.StationID)
			if err != nil {
				// If it's not in the static databse, it's an outpost.
				sta, err = e.xmlAPI.OutpostForID(o.StationID)
				if err != nil {
					// Make a dummy station.
					sta = &types.Station{
						Name: fmt.Sprintf("Unknown Station (ID %d)", o.StationID),
						ID:   o.StationID,
					}
				}
			}
			stationCache[o.StationID] = sta
		}
		oTime, _ := time.Parse("2006-01-02", o.ExpirationDate)
		newOrder := types.Order{
			Type:       t,
			Item:       item,
			Quantity:   o.QuantityAvailable,
			Station:    stationCache[o.StationID],
			Price:      o.Price,
			Expiration: oTime,
		}
		if t == types.Buy {
			// Set the fields specific to buy orders.
			newOrder.MinQuantity = o.MinimumVolume
			switch o.Range {
			case 32767, 65535:
				newOrder.JumpRange = types.BuyRegion
			case -1:
				newOrder.JumpRange = types.BuyStation
			case 0:
				newOrder.JumpRange = types.BuySystem
			default:
				newOrder.JumpRange = types.BuyNumberJumps
				newOrder.NumJumps = o.Range
			}
		}
		results = append(results, newOrder)
	}
	return results
}

func (e *eveCentral) OrdersForItem(item *types.Item, location string, orderType types.OrderType) (*[]types.Order, error) {
	var (
		system *types.SolarSystem
		region *types.Region
		err    error
	)
	system, err = e.db.SolarSystemForName(location)
	if err != nil {
		// Not a system or unable to look up. Try region.
		region, err = e.db.RegionForName(location)
		if err != nil {
			// Still can't find it. Return an error.
			return nil, err
		}
	}
	query := url.Values{}
	if region != nil {
		query.Set("regionlimit", fmt.Sprintf("%d", region.ID))
	} else {
		query.Set("usesystem", fmt.Sprintf("%d", system.ID))
	}
	query.Set("typeid", fmt.Sprintf("%d", item.ID))
	e.endpoint.RawQuery = query.Encode()
	orderXML, found := e.respCache.Get(query.Encode())
	if !found {
		var err error
		orderXML, err = e.getURL(e.endpoint.String())
		if err != nil {
			return nil, err
		}
		// EVE-Central doesn't specify a caching time to use, so we're picking
		// five minutes at random.
		e.respCache.Put(query.Encode(), orderXML, time.Now().Add(5*time.Minute))
	}
	orders := &quicklook{}

	err = xml.Unmarshal(orderXML, orders)
	if err != nil {
		return nil, err
	}

	// Convert returned XML struct into what we present to rest of library.
	results := []types.Order{}
	switch orderType {
	case types.AllOrders:
		// The order here matters, if only because it's the order that the
		// orders are presented by EVE Central and therefore the order in which
		// the test cases expect results.
		results = append(results, e.processOrders(orders, item, types.Sell)...)
		results = append(results, e.processOrders(orders, item, types.Buy)...)
	default:
		results = e.processOrders(orders, item, orderType)
	}
	return &results, nil
}

func (e *eveCentral) BuyInStation(item *types.Item, location *types.Station) (*[]types.Order, error) {
	system, err := e.db.SolarSystemForID(location.SystemID)
	if err != nil {
		return nil, err
	}
	regionalOrders, err := e.OrdersForItem(item, system.Region, types.Buy)
	if err != nil {
		return nil, err
	}
	orders := []types.Order{}
	for _, o := range *regionalOrders {
		switch o.JumpRange {
		case types.BuyRegion:
			orders = append(orders, o)
		case types.BuyNumberJumps:
			numJumps, err := e.router.NumJumpsID(o.Station.SystemID, location.SystemID)
			if err != nil {
				return nil, err
			}
			if numJumps <= o.NumJumps {
				orders = append(orders, o)
			}
		case types.BuySystem:
			if o.Station.SystemID == location.SystemID {
				orders = append(orders, o)
			}
		case types.BuyStation:
			if o.Station.ID == location.ID {
				orders = append(orders, o)
			}
		}
	}

	return &orders, nil
}

func (e *eveCentral) OrdersInStation(item *types.Item, location *types.Station) (*[]types.Order, error) {
	orders, err := e.BuyInStation(item, location)
	if err != nil {
		return nil, err
	}
	// Get the sell orders for the entire system, then append the ones for this
	// station to the returned array.
	orderSystem, err := e.db.SolarSystemForID(location.SystemID)
	if err != nil {
		return nil, err
	}
	sellInSystem, err := e.OrdersForItem(item, orderSystem.Name, types.Sell)
	if err != nil {
		return nil, err
	}
	for _, o := range *sellInSystem {
		if o.Station.ID == location.ID {
			*orders = append(*orders, o)
		}
	}
	return orders, nil
}

func (e *eveCentral) Close() error {
	return nil
}
