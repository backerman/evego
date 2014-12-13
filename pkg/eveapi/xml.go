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

// Package eveapi is the public interface for accessing the
// EVE APIs (XML, CREST, or whatever.)
package eveapi

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/types"
)

var (
	cacheExpiry = time.Time{}
	outposts    = make(map[int]*types.Station)
)

const (
	iso8601            = "2006-01-02 15:04:05"
	conqerableStations = "/eve/ConquerableStationList.xml.aspx"
)

type xmlAPI struct {
	// Endpoint URL to access.
	url  *url.URL
	http http.Client
	db   dbaccess.EveDatabase
}

// XMLAPI returns an EveAPI that accesses the EVE Online XML API.
func XMLAPI(serviceURL string, staticDB dbaccess.EveDatabase) EveAPI {
	endpoint, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatalf("Unable to process endpoint URL: %v", err)
	}
	return &xmlAPI{url: endpoint, db: staticDB}
}

// XML headers for unmarshalling

type apiResponse struct {
	CurrentTime string    `xml:"currentTime"`
	Outposts    []outpost `xml:"result>rowset>row"`
	CachedUntil string    `xml:"cachedUntil"`
}

type outpost struct {
	Name            string `xml:"stationName,attr"`
	ID              int    `xml:"stationID,attr"`
	TypeID          int    `xml:"stationTypeID,attr"`
	SolarSystemID   int    `xml:"solarSystemID,attr"`
	CorporationID   int    `xml:"corporationID,attr"`
	CorporationName string `xml:"corporationName,attr"`
}

func expirationTime(currentTime, cachedUntil string) time.Time {
	// This system's time *should* be identical to the server's,
	// but we're not taking any chances.
	current, err := time.Parse(iso8601, currentTime)
	if err != nil {
		// If we can't parse, return the datum.
		log.Printf("Unable to parse time from server: %v", currentTime)
		return time.Time{}
	}
	until, err := time.Parse(iso8601, cachedUntil)
	if err != nil {
		log.Printf("Unable to parse time from server: %v", cachedUntil)
		return time.Time{}
	}
	diff := until.Sub(current)
	return time.Now().Add(diff)
}

func (x *xmlAPI) get(endpoint string, params ...url.Values) ([]byte, error) {
	// Make a copy of our base URL and modify as appropriate for this call.
	callURL := *x.url
	callURL.Path = endpoint
	var (
		req *http.Request
		err error
	)
	if params == nil {
		req, err = http.NewRequest("GET", callURL.String(), nil)
	} else {
		req, err = http.NewRequest("POST", callURL.String(),
			strings.NewReader(params[0].Encode()))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "evego (https://github.com/backerman/evego)")
	if params != nil {
		// We're doing a POST so need to set the content type.
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := x.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (x *xmlAPI) OutpostForID(id int) (*types.Station, error) {
	if time.Now().After(cacheExpiry) {
		// The cache has expired or has not yet been populated.
		// FIXME Only one goroutine should update cache. Mutex? Channel?
		newOutposts := make(map[int]*types.Station)
		xmlBytes, err := x.get(conqerableStations)
		if err != nil {
			// FIXME some sort of throttling required
			return nil, err
		}
		var response apiResponse
		xml.Unmarshal(xmlBytes, &response)
		cacheExpiry = expirationTime(response.CurrentTime, response.CachedUntil)
		for i := range response.Outposts {
			o := response.Outposts[i]
			stn := types.Station{
				Name:     o.Name,
				ID:       o.ID,
				SystemID: o.SolarSystemID,
				// Delay constellation/region lookup until queried.
			}
			newOutposts[o.ID] = &stn
		}
		outposts = newOutposts
	}
	stn, exists := outposts[id]
	if !exists {
		return nil, fmt.Errorf("Station ID %d not found.", id)
	}
	if stn.ConstellationID == 0 {
		system, err := x.db.SolarSystemForID(stn.SystemID)
		if err != nil {
			return nil, err
		}
		stn.ConstellationID = system.ConstellationID
		stn.RegionID = system.RegionID
	}
	return stn, nil
}

func (x *xmlAPI) Close() error {
	return nil
}
