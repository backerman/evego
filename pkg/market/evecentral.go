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

	"github.com/backerman/evego"
)

type eveCentral struct {
	db        evego.Database
	router    evego.Router
	xmlAPI    evego.XMLAPI
	endpoint  *url.URL
	http      http.Client
	respCache evego.Cache
}

// EveCentral returns an interface to the EVE-Central API.
// It takes as input an EveDatabase object and an HTTP endpoint;
// the latter should be http://api.eve-central.com/api/quicklook
// for the production EVE-Central instance.
func EveCentral(db evego.Database, router evego.Router, xmlAPI evego.XMLAPI, endpoint string,
	aCache evego.Cache) evego.Market {
	epURL, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("Invalid URL %v passed for Eve-Central endpoint: %v", endpoint, err)
	}
	ec := eveCentral{db: db, router: router, endpoint: epURL, xmlAPI: xmlAPI, respCache: aCache}
	return &ec
}

func (e *eveCentral) getURL(u string) ([]byte, error) {
	// Start by checking cache.
	body, found := e.respCache.Get(u)
	if found {
		return body, nil
	}
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
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// EVE-Central doesn't specify a caching time to use, so we're picking
		// ten minutes at random.
		e.respCache.Put(u, body, time.Now().Add(10*time.Minute))
	}
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

func (e *eveCentral) processOrders(data *quicklook, item *evego.Item, t evego.OrderType) []evego.Order {
	var toProcess *[]order
	// Set up a temporary cache so that we only get each station's object once.
	stationCache := make(map[int]*evego.Station)
	switch t {
	case evego.Buy:
		toProcess = &data.BuyOrders
	case evego.Sell:
		toProcess = &data.SellOrders
	}
	results := []evego.Order{}
	for _, o := range *toProcess {
		if stationCache[o.StationID] == nil {
			sta, err := e.db.StationForID(o.StationID)
			if err != nil {
				// If it's not in the static databse, it's an outpost.
				sta, err = e.xmlAPI.OutpostForID(o.StationID)
				if err != nil {
					// Make a dummy station.
					sta = &evego.Station{
						Name: fmt.Sprintf("Unknown Station (ID %d)", o.StationID),
						ID:   o.StationID,
					}
				}
			}
			stationCache[o.StationID] = sta
		}
		oTime, _ := time.Parse("2006-01-02", o.ExpirationDate)
		newOrder := evego.Order{
			Type:       t,
			Item:       item,
			Quantity:   o.QuantityAvailable,
			Station:    stationCache[o.StationID],
			Price:      o.Price,
			Expiration: oTime,
		}
		if t == evego.Buy {
			// Set the fields specific to buy orders.
			newOrder.MinQuantity = o.MinimumVolume
			switch o.Range {
			case 32767, 65535:
				newOrder.JumpRange = evego.BuyRegion
			case -1:
				newOrder.JumpRange = evego.BuyStation
			case 0:
				newOrder.JumpRange = evego.BuySystem
			default:
				newOrder.JumpRange = evego.BuyNumberJumps
				newOrder.NumJumps = o.Range
			}
		}
		results = append(results, newOrder)
	}
	return results
}

func (e *eveCentral) OrdersForItem(item *evego.Item, location string, orderType evego.OrderType) (*[]evego.Order, error) {
	var (
		system *evego.SolarSystem
		region *evego.Region
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
	orderXML, err := e.getURL(e.endpoint.String())
	if err != nil {
		return nil, err
	}

	// We received a quicklook XML document from EVE-Central. Unmarshal it.
	orders := &quicklook{}
	err = xml.Unmarshal(orderXML, orders)
	if err != nil {
		return nil, err
	}

	// Convert returned XML struct into what we present to rest of library.
	results := []evego.Order{}
	switch orderType {
	case evego.AllOrders:
		// The order here matters, if only because it's the order that the
		// orders are presented by EVE Central and therefore the order in which
		// the test cases expect results.
		results = append(results, e.processOrders(orders, item, evego.Sell)...)
		results = append(results, e.processOrders(orders, item, evego.Buy)...)
	default:
		results = e.processOrders(orders, item, orderType)
	}
	return &results, nil
}

func (e *eveCentral) BuyInStation(item *evego.Item, location *evego.Station) (*[]evego.Order, error) {
	system, err := e.db.SolarSystemForID(location.SystemID)
	if err != nil {
		return nil, err
	}
	regionalOrders, err := e.OrdersForItem(item, system.Region, evego.Buy)
	if err != nil {
		return nil, err
	}
	orders := []evego.Order{}
	for _, o := range *regionalOrders {
		switch o.JumpRange {
		case evego.BuyRegion:
			orders = append(orders, o)
		case evego.BuyNumberJumps:
			numJumps, err := e.router.NumJumpsID(o.Station.SystemID, location.SystemID)
			if err != nil {
				return nil, err
			}
			if numJumps <= o.NumJumps {
				orders = append(orders, o)
			}
		case evego.BuySystem:
			if o.Station.SystemID == location.SystemID {
				orders = append(orders, o)
			}
		case evego.BuyStation:
			if o.Station.ID == location.ID {
				orders = append(orders, o)
			}
		}
	}

	return &orders, nil
}

func (e *eveCentral) OrdersInStation(item *evego.Item, location *evego.Station) (*[]evego.Order, error) {
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
	sellInSystem, err := e.OrdersForItem(item, orderSystem.Name, evego.Sell)
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
